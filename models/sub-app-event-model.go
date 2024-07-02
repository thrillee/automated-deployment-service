package models

import (
	"context"

	"github.com/thrillee/automated-deployment-service/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	SUBSCRIBER_APP_PENDING   = "SUBSCRIBER_APP_PENDING"
	SUBSCRIBER_APP_DEPLOYING = "SUBSCRIBER_APP_DEPLOYING"
	SUBSCRIBER_APP_COMPLETD  = "SUBSCRIBER_APP_COMPLETD"
	SUBSCRIBER_APP_FAILED    = "SUBSCRIBER_APP_FAILED"
)

type SubscriberAppEvent struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AppSubscriberID primitive.ObjectID `bson:"app_subscriber_id" json:"app_subscriber_id"`
	AppDeployStepID primitive.ObjectID `bson:"app_deploy_step_id" json:"app_deploy_step_id"`
	Step            int                `bson:"step" json:"step"`
	Status          string             `bson:"status" json:"status"`
}

func (ase *SubscriberAppEvent) Insert(ctx context.Context, db *db.MongoDB) error {
	collection := db.GetCollection("subscriber_app_events")
	ase.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(ctx, ase)
	return err
}

func (ase *SubscriberAppEvent) Update(ctx context.Context, db *db.MongoDB) error {
	collection := db.GetCollection("subscriber_app_events")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": ase.ID}, bson.M{"$set": ase})
	return err
}

func ListAppSubscriberEventsByAppAndSubscriber(ctx context.Context, db *db.MongoDB, appID, subscriberID primitive.ObjectID) ([]SubscriberAppEvent, error) {
	collection := db.GetCollection("subscriber_app_events")
	pipeline := mongo.Pipeline{
		{{"$lookup", bson.M{
			"from":         "subscriber_apps",
			"localField":   "subscriber_app_id",
			"foreignField": "_id",
			"as":           "subscriber_app",
		}}},
		{{"$unwind", "$subscriber_app"}},
		{{"$match", bson.M{
			"subscriber_app.app_id":        appID,
			"subscriber_app.subscriber_id": subscriberID,
		}}},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var results []SubscriberAppEvent
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
