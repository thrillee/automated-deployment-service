package models

import (
	"context"
	"time"

	"github.com/thrillee/automated-deployment-service/internals/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	SUBSCRIBER_APP_EVENT_PENDING   = "SUBSCRIBER_APP_EVENT_PENDING"
	SUBSCRIBER_APP_EVENT_DEPLOYING = "SUBSCRIBER_APP_EVENT_DEPLOYING"
	SUBSCRIBER_APP_EVENT_COMPLETD  = "SUBSCRIBER_APP_EVENT_COMPLETD"
	SUBSCRIBER_APP_EVENT_FAILED    = "SUBSCRIBER_APP_EVENT_FAILED"
)

type SubscriberAppEvent struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AppSubscriberID primitive.ObjectID `bson:"app_subscriber_id" json:"app_subscriber_id"`
	AppDeployStepID primitive.ObjectID `bson:"app_deploy_step_id" json:"app_deploy_step_id"`
	Step            int                `bson:"step" json:"step"`
	Status          string             `bson:"status" json:"status"`
	Modified        time.Time          `bson:"modified" json:"modified"`
	Created         time.Time          `bson:"created" json:"created"`
}

func BulkInsertSubscriberAppEvent(ctx context.Context, db *db.MongoDB, events []SubscriberAppEvent) error {
	collection := db.GetCollection("subscriber_app_events")
	docs := make([]interface{}, len(events))
	for i, e := range events {
		e.ID = primitive.NewObjectID()
		e.Modified = time.Now()
		e.Created = time.Now()
		docs[i] = e
	}
	_, err := collection.InsertMany(ctx, docs)

	return err
}

func (ase *SubscriberAppEvent) Insert(ctx context.Context, db *db.MongoDB) error {
	collection := db.GetCollection("subscriber_app_events")
	ase.ID = primitive.NewObjectID()
	ase.Created = time.Now()
	ase.Modified = time.Now()
	_, err := collection.InsertOne(ctx, ase)
	return err
}

func (ase *SubscriberAppEvent) Update(ctx context.Context, db *db.MongoDB) error {
	collection := db.GetCollection("subscriber_app_events")
	ase.Modified = time.Now()
	_, err := collection.UpdateOne(ctx, bson.M{"_id": ase.ID}, bson.M{"$set": ase})
	return err
}

func ListAppSubscriberEventsByStatus(ctx context.Context, db *db.MongoDB, subAppId primitive.ObjectID, status string) ([]SubscriberAppEvent, error) {
	collection := db.GetCollection("subscriber_app_events")
	cursor, err := collection.Find(
		ctx,
		bson.M{"status": status, "app_subscriber_id": subAppId},
		options.Find().SetSort(bson.D{{Key: "step", Value: 1}}))

	var results []SubscriberAppEvent
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
