package models

import (
	"context"
	"time"

	"github.com/thrillee/automated-deployment-service/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subscriber struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title" validate:"required"`
	Description string             `bson:"description" json:"description" validate:"required"`
	Reference   string             `bson:"reference" json:"reference"`
	CallerRef   string             `bson:"callerref" json:"caller_ref"`
	Modified    time.Time          `bson:"modified" json:"modified"`
	Created     time.Time          `bson:"created" json:"created"`
}

func (s *Subscriber) FindSubscriberByReference(ctx context.Context, db *db.MongoDB, reference string) error {
	collection := db.GetCollection("subscribers")
	return collection.FindOne(ctx, bson.M{"reference": reference}).Decode(&s)
}

func (s *Subscriber) FindSubscriberByCallerRef(ctx context.Context, db *db.MongoDB, callerRef string) error {
	collection := db.GetCollection("subscribers")
	return collection.FindOne(ctx, bson.M{"callerref": callerRef}).Decode(&s)
}

func (s *Subscriber) Insert(ctx context.Context, db *db.MongoDB) error {
	s.Created = time.Now()
	s.Modified = time.Now()
	collection := db.GetCollection("subscribers")
	s.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(ctx, s)
	return err
}

func (s *Subscriber) Update(ctx context.Context, db *db.MongoDB) error {
	s.Modified = time.Now()
	collection := db.GetCollection("subscribers")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": s.ID}, bson.M{"$set": s})
	return err
}

func (s *Subscriber) Delete(ctx context.Context, db *db.MongoDB) error {
	collection := db.GetCollection("subscribers")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": s.ID})
	return err
}
