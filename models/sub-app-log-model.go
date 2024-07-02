package models

import (
	"context"

	"github.com/thrillee/automated-deployment-service/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubscriberAppLog struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AppSubscriberID primitive.ObjectID `bson:"app_subscriber_id" json:"app_subscriber_id"`
	AppDeployStepID primitive.ObjectID `bson:"app_deploy_step_id" json:"app_deploy_step_id"`
	NKey            string             `bson:"nkey" json:"n_key"`
	NValue          string             `bson:"nvalue" json:"n_value"`
}

func ListAppSubscriberLogsByAppAndSubscriber(ctx context.Context, db *db.MongoDB, appID, subscriberID primitive.ObjectID) ([]SubscriberAppLog, error) {
	collection := db.GetCollection("subscriber_app_logs")
	pipeline := mongo.Pipeline{
		{{"$lookup", bson.M{
			"from":         "app_subscribers",
			"localField":   "app_subscriber_id",
			"foreignField": "_id",
			"as":           "app_subscriber",
		}}},
		{{"$unwind", "$app_subscriber"}},
		{{"$match", bson.M{
			"app_subscriber.app_id":        appID,
			"app_subscriber.subscriber_id": subscriberID,
		}}},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var results []SubscriberAppLog
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func FilterAppSubscriberLogsByAppSubscriberAndNKey(ctx context.Context, db *db.MongoDB, appID, subscriberID primitive.ObjectID, nkey string) ([]SubscriberAppLog, error) {
	collection := db.GetCollection("subscriber_app_logs")
	pipeline := mongo.Pipeline{
		{{"$lookup", bson.M{
			"from":         "app_subscribers",
			"localField":   "app_subscriber_id",
			"foreignField": "_id",
			"as":           "app_subscriber",
		}}},
		{{"$unwind", "$app_subscriber"}},
		{{"$match", bson.M{
			"app_subscriber.app_id":        appID,
			"app_subscriber.subscriber_id": subscriberID,
			"nkey":                         nkey,
		}}},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var results []SubscriberAppLog
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (asl *SubscriberAppLog) Insert(ctx context.Context, db *db.MongoDB) error {
	collection := db.GetCollection("subscriber_app_logs")
	_, err := collection.InsertOne(ctx, asl)
	return err
}

func (asl *SubscriberAppLog) Update(ctx context.Context, db *db.MongoDB) error {
	collection := db.GetCollection("subscriber_app_logs")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": asl.ID}, bson.M{"$set": asl})
	return err
}

func (asl *SubscriberAppLog) Delete(ctx context.Context, db *db.MongoDB) error {
	collection := db.GetCollection("subscriber_app_logs")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": asl.ID})
	return err
}
