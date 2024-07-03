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
	SUBSCRIBER_APP_PENDING   = "SUBSCRIBER_APP_PENDING"
	SUBSCRIBER_APP_DEPLOYING = "SUBSCRIBER_APP_DEPLOYING"
	SUBSCRIBER_APP_COMPLETD  = "SUBSCRIBER_APP_COMPLETD"
	SUBSCRIBER_APP_FAILED    = "SUBSCRIBER_APP_FAILED"
)

var sub_app_collection = "subscriber_app"

type SubscriberApp struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AppID          primitive.ObjectID `bson:"app_id" json:"app_id"`
	SubscriberID   primitive.ObjectID `bson:"subscriber_id" json:"subscriber_id"`
	Status         string             `bson:"status" json:"status"`
	DeployPriority int                `bson:"deploy_priority" json:"deploy_priority"`
	Modified       time.Time          `bson:"modified" json:"modified"`
	Created        time.Time          `bson:"created" json:"created"`
}

func ListSubAppsByStatus(ctx context.Context, db *db.MongoDB, status string) ([]SubscriberApp, error) {
	collection := db.GetCollection(sub_app_collection)
	cursor, err := collection.Find(ctx, bson.M{"Status": status}, options.Find().SetSort(bson.D{{Key: "deploy_priority", Value: 1}}))
	if err != nil {
		return nil, err
	}

	var subApps []SubscriberApp
	if err := cursor.All(ctx, &subApps); err != nil {
		return nil, err
	}
	return subApps, nil
}

func ListAppSubscribersBySubscriberID(ctx context.Context, db *db.MongoDB, subscriberID primitive.ObjectID) ([]SubscriberApp, error) {
	collection := db.GetCollection(sub_app_collection)
	cursor, err := collection.Find(ctx, bson.M{"subscriber_id": subscriberID})
	if err != nil {
		return nil, err
	}
	var results []SubscriberApp
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (as *SubscriberApp) Insert(ctx context.Context, db *db.MongoDB) error {
	as.Created = time.Now()
	as.Modified = time.Now()
	as.ID = primitive.NewObjectID()
	collection := db.GetCollection(sub_app_collection)
	_, err := collection.InsertOne(ctx, as)
	return err
}

func (as *SubscriberApp) Update(ctx context.Context, db *db.MongoDB) error {
	as.Modified = time.Now()
	collection := db.GetCollection(sub_app_collection)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": as.ID}, bson.M{"$set": as})
	return err
}

func (as *SubscriberApp) Delete(ctx context.Context, db *db.MongoDB) error {
	collection := db.GetCollection(sub_app_collection)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": as.ID})
	return err
}
