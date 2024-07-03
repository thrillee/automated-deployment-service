package models

import (
	"context"

	"github.com/thrillee/automated-deployment-service/internals/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name                 string             `bson:"name" json:"name"`
	Description          string             `bson:"description" json:"description"`
	Reference            string             `bson:"reference" json:"reference"`
	GithubRepo           string             `bson:"github_repo" json:"github_repo"`
	ContainerRegistryURL string             `bson:"container_registry_url" json:"container_registry_url"`
	Processor            string             `bson:"processor" json:"processor"`
	DeployPriority       int                `bson:"deploy_priority" json:"deploy_priority"`
}

func (a *App) FindAppById(ctx context.Context, db *db.MongoDB, id primitive.ObjectID) error {
	collection := db.GetCollection("apps")
	return collection.FindOne(ctx, bson.M{"_id": id}).Decode(&a)
}

func (a *App) FindAppByReference(ctx context.Context, db *db.MongoDB, reference string) error {
	collection := db.GetCollection("apps")
	return collection.FindOne(ctx, bson.M{"reference": reference}).Decode(&a)
}

var app_deploy_steps = "app_deploy_steps"

type AppDeployStep struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	AppID       primitive.ObjectID `bson:"app_id"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	NKey        string             `bson:"nkey"`
	Step        int                `bson:"step"`
}

func (a *AppDeployStep) FindById(ctx context.Context, db *db.MongoDB, id primitive.ObjectID) error {
	collection := db.GetCollection(app_deploy_steps)
	return collection.FindOne(ctx, bson.M{"_id": id}).Decode(&a)
}

func ListAppDeployStepsByAppID(ctx context.Context, db *db.MongoDB, appID primitive.ObjectID) ([]AppDeployStep, error) {
	collection := db.GetCollection(app_deploy_steps)
	cursor, err := collection.Find(ctx, bson.M{"app_id": appID}, options.Find().SetSort(bson.D{{Key: "step", Value: 1}}))
	if err != nil {
		return nil, err
	}
	var results []AppDeployStep
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
