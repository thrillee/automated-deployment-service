package models

import (
	"context"

	"github.com/thrillee/automated-deployment-service/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty"`
	Name                 string             `bson:"name"`
	Description          string             `bson:"description"`
	Reference            string             `bson:"reference"`
	GithubRepo           string             `bson:"github_repo"`
	ContainerRegistryURL string             `bson:"container_registry_url"`
}

func (a *App) FindAppByReference(ctx context.Context, db *db.MongoDB, reference string) error {
	collection := db.GetCollection("apps")
	return collection.FindOne(ctx, bson.M{"reference": reference}).Decode(&a)
}

type AppDeployStep struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	AppID       primitive.ObjectID `bson:"app_id"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	NKey        string             `bson:"nkey"`
	Step        int                `bson:"step"`
}

func ListAppDeployStepsByAppID(ctx context.Context, db *db.MongoDB, appID primitive.ObjectID) ([]AppDeployStep, error) {
	collection := db.GetCollection("app_deploy_steps")
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
