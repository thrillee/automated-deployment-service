package processor

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/thrillee/automated-deployment-service/internals/models"
)

type handlerPayload struct {
	Subscriber    models.Subscriber
	appDeployStep models.AppDeployStep
	app           models.App
}

func (d *DeploymentProcessor) PickDeployables() {
	batch_id := uuid.New().String()
	ctx := context.Background()
	ctx = context.WithValue(ctx, processs_batch_id_key, batch_id)

	pendingSubApps, err := models.ListSubAppsByStatus(ctx, d.db, models.SUBSCRIBER_APP_DEPLOYING)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"process-batch-id": getProcessBatchId(ctx),
			"action":           "Handle Sub App Event: List Subscriber App",
		}).Error(err)
		return
	}

	for _, subApp := range pendingSubApps {
		go d.handleSubAppEvents(ctx, subApp)
	}
}

func (d *DeploymentProcessor) handleSubAppEvents(ctx context.Context, subApp models.SubscriberApp) {
	events, err := models.ListAppSubscriberEventsByStatus(ctx, d.db, subApp.ID, models.SUBSCRIBER_APP_PENDING)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"process-batch-id": getProcessBatchId(ctx),
			"action":           "Handle Sub App Event: List Subscriber App Events",
			"sub-app":          subApp,
		}).Error(err)
		return
	}

	for _, e := range events {
		d.handleAppEvents(ctx, subApp, e)
	}
}

func (d *DeploymentProcessor) handleAppEvents(ctx context.Context, subApp models.SubscriberApp, subAppEvent models.SubscriberAppEvent) {
	subscriber := models.Subscriber{}
	err := subscriber.FindById(ctx, d.db, subApp.SubscriberID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"process-batch-id": getProcessBatchId(ctx),
			"action":           "Handle Sub-App-Event: Find Subscriber",
			"sub-app":          subApp,
			"sub-app-event":    subAppEvent,
		}).Error(err)
		return
	}

	appDeployStep := models.AppDeployStep{}
	err = appDeployStep.FindById(ctx, d.db, subAppEvent.AppDeployStepID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"process-batch-id": getProcessBatchId(ctx),
			"action":           "Handle Sub-App-Event: Fetch App Deployment Steps",
			"subscriber":       subscriber,
			"sub-app":          subApp,
			"sub-app-event":    subAppEvent,
		}).Error(err)
		return
	}

	app := models.App{}
	err = app.FindAppById(ctx, d.db, appDeployStep.AppID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"process-batch-id": getProcessBatchId(ctx),
			"action":           "Handle Sub-App-Event: Fetch App",
			"subscriber":       subscriber,
			"sub-app":          subApp,
			"sub-app-event":    subAppEvent,
			"app-deploy-step":  appDeployStep,
		}).Error(err)
		return
	}

	handlerFunc, err := d.processorFactory.GetHandler(appDeployStep.NKey)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"process-batch-id": getProcessBatchId(ctx),
			"action":           "Handle Sub-App-Event: Get Handler",
			"subscriber":       subscriber,
			"sub-app":          subApp,
			"sub-app-event":    subAppEvent,
			"app-deploy-step":  appDeployStep,
			"app-deploy-key":   appDeployStep.NKey,
		}).Error(err)
		return
	}

	hp := handlerPayload{
		app:           app,
		Subscriber:    subscriber,
		appDeployStep: appDeployStep,
	}

	result := handlerFunc(ctx, hp)
	if !result.Success {
		subAppEvent.Status = models.SUBSCRIBER_APP_EVENT_PENDING
		subAppEvent.Update(ctx, d.db)
	} else {
		subAppEvent.Status = models.SUBSCRIBER_APP_EVENT_COMPLETD
		subAppEvent.Update(ctx, d.db)
	}

	logrus.WithFields(logrus.Fields{
		"process-batch-id":          getProcessBatchId(ctx),
		"action":                    "Handle Sub-App-Event: Completed",
		"subscriber":                subscriber,
		"sub-app":                   subApp,
		"sub-app-event":             subAppEvent,
		"app-deploy-step":           appDeployStep,
		"app-deploy-key":            appDeployStep.NKey,
		"app-deploy-handler":        handlerFunc,
		"app-deploy-handler-input":  hp,
		"app-deploy-handler-result": result,
	}).Info("Completed Handle App Event")
}
