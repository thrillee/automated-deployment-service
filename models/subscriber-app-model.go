package models

import (
	"context"

	"github.com/thrillee/automated-deployment-service/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubscriberApp struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AppID        primitive.ObjectID `bson:"app_id" json:"app_id"`
	SubscriberID primitive.ObjectID `bson:"subscriber_id" json:"subscriber_id"`
	Status       string             `bson:"status" json:"status"`
}

func ListAppSubscribersBySubscriberID(ctx context.Context, db *db.MongoDB, subscriberID primitive.ObjectID) ([]SubscriberApp, error) {
	collection := db.GetCollection("subscriber_app")
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

func FilterAppSubscribersBySubscriberIDAndAppReference(ctx context.Context, db *db.MongoDB, subscriberID primitive.ObjectID, appReference string) ([]SubscriberApp, error) {
	collection := db.GetCollection("subscriber_apps")
	pipeline := mongo.Pipeline{
		{{"$lookup", bson.M{
			"from":         "apps",
			"localField":   "app_id",
			"foreignField": "_id",
			"as":           "app",
		}}},
		{{"$unwind", "$app"}},
		{{"$match", bson.M{
			"subscriber_id": subscriberID,
			"app.reference": appReference,
		}}},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
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
	collection := db.GetCollection("subscriber_apps")
	_, err := collection.InsertOne(ctx, as)
	return err
}

func (as *SubscriberApp) Update(ctx context.Context, db *db.MongoDB) error {
	collection := db.GetCollection("subscriber_apps")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": as.ID}, bson.M{"$set": as})
	return err
}

func (as *SubscriberApp) Delete(ctx context.Context, db *db.MongoDB) error {
	collection := db.GetCollection("subscriber_apps")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": as.ID})
	return err
}
