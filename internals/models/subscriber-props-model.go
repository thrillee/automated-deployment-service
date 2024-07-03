package models

import (
	"context"
	"time"

	"github.com/thrillee/automated-deployment-service/internals/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubscriberProp struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SubscriberID primitive.ObjectID `bson:"subscriber_id" json:"subscriber_id"`
	NKey         string             `bson:"nkey" json:"n_key"`
	NValue       string             `bson:"nvalue" json:"n_value"`
	Modified     time.Time          `bson:"modified" json:"modified"`
	Created      time.Time          `bson:"created" json:"created"`
}

func (sp *SubscriberProp) FindSubscrierPropByKey(ctx context.Context, db *db.MongoDB, subscriberID primitive.ObjectID, nkey string) error {
	collection := db.GetCollection("subscriber_props")
	return collection.FindOne(ctx, bson.M{"subscriber_id": subscriberID, "nkey": nkey}).Decode(sp)
}

func (sp *SubscriberProp) Insert(ctx context.Context, db *db.MongoDB) error {
	sp.Created = time.Now()
	sp.Modified = time.Now()
	collection := db.GetCollection("subscriber_props")
	sp.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(ctx, sp)
	return err
}

func (sp *SubscriberProp) Update(ctx context.Context, db *db.MongoDB) error {
	sp.Modified = time.Now()
	collection := db.GetCollection("subscriber_props")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": sp.ID}, bson.M{"$set": sp})
	return err
}

func (sp *SubscriberProp) Delete(ctx context.Context, db *db.MongoDB) error {
	collection := db.GetCollection("subscriber_props")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": sp.ID})
	return err
}

func ListSubscriberPropsBySubscriberID(ctx context.Context, db *db.MongoDB, subscriberID primitive.ObjectID) ([]SubscriberProp, error) {
	collection := db.GetCollection("subscriber_props")
	cursor, err := collection.Find(ctx, bson.M{"subscriber_id": subscriberID})
	if err != nil {
		return nil, err
	}
	var results []SubscriberProp
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
