package controllers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/thrillee/automated-deployment-service/internals/models"
)

type NewSubscriberPayload struct {
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Reference   string   `json:"reference" validate:"required"`
	Subdomain   string   `json:"subdomain" validate:"required"`
	Apps        []string `json:"apps" validate:"required"`
}

func (c *Controller) AddSubscriber(ctx context.Context, data NewSubscriberPayload) ControllerResponse {
	if validationErr := validate.Struct(&data); validationErr != nil {
		return ControllerResponse{
			Success: false,
			Message: validationErr.Error(),
		}
	}

	newSub := models.Subscriber{
		Title:       data.Title,
		Description: data.Description,
		Reference:   uuid.New().String(),
		CallerRef:   data.Reference,
	}
	newSub.Insert(ctx, c.db)

	props := models.SubscriberProp{
		SubscriberID: newSub.ID,
		NKey:         "sub_domain",
		NValue:       data.Subdomain,
	}
	props.Insert(ctx, c.db)

	subAppEvents := []models.SubscriberAppEvent{}
	for _, appRef := range data.Apps {
		app := models.App{}
		err := app.FindAppByReference(ctx, c.db, appRef)
		if err != nil {
			return ControllerResponse{
				Success: false,
				Message: fmt.Sprintf("Failed to find Ref: %s Error: %v", appRef, err),
			}
		}

		deployStep, err := models.ListAppDeployStepsByAppID(ctx, c.db, app.ID)
		if err != nil {
			return ControllerResponse{
				Success: false,
				Message: fmt.Sprintf("Failed to steps: %s Error: %v", appRef, err),
			}
		}

		subApp := models.SubscriberApp{
			AppID:          app.ID,
			SubscriberID:   newSub.ID,
			DeployPriority: app.DeployPriority,
			Status:         models.SUBSCRIBER_APP_DEPLOYING,
		}
		subApp.Insert(ctx, c.db)

		for _, d := range deployStep {
			logrus.WithFields(logrus.Fields{
				"App":        data.Title,
				"Step Title": d.Title,
				"Step":       d.Step,
			}).Info("Step")

			ads := models.SubscriberAppEvent{
				AppSubscriberID: subApp.ID,
				AppDeployStepID: d.ID,
				Step:            d.Step,
				Status:          models.SUBSCRIBER_APP_EVENT_PENDING,
			}
			subAppEvents = append(subAppEvents, ads)
		}

	}

	err := models.BulkInsertSubscriberAppEvent(ctx, c.db, subAppEvents)
	if err != nil {
		return ControllerResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to log deployment steps: %v", err),
		}
	}

	logrus.WithFields(logrus.Fields{
		"app-deployment-steps": len(subAppEvents),
		"app_count":            len(data.Apps),
	}).Info("Deployment steps count")

	subSteps := []string{}

	result := map[string]interface{}{}
	result["subscriberReference"] = newSub.Reference
	result["appSub"] = subSteps

	return ControllerResponse{
		Success: true,
		Message: "Subscriber Added Successfully",
		Result:  result,
	}
}
