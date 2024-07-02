package db

import (
	"context"
	"log"
	"time"

	"github.com/thrillee/automated-deployment-service/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var client *mongo.Client

var db *MongoDB

type MongoDB struct {
	client *mongo.Client
	ctx    context.Context
	dbName string
}

func GetClient() (*mongo.Client, error) {
	var err error
	if db == nil {
		db, err = New()
		if err != nil {
			return nil, err
		}
	}

	return db.client, nil
}

func New() (*MongoDB, error) {
	if db == nil {
		dbConfig := common.GetDBConfig()
		client, ctx, err := newMongoDBConnection(dbConfig)
		if err != nil {
			return nil, err
		}
		db = &MongoDB{
			client: client,
			ctx:    ctx,
			dbName: db.dbName,
		}
	}

	return db, nil
}

func (db *MongoDB) GetCollection(collectionName string) *mongo.Collection {
	collection := db.client.Database("ads").Collection(collectionName)
	return collection
}

func (db *MongoDB) PingDB(collectionName string) {
	err := db.client.Ping(db.ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Pinged your deployment. You successfully connected to MongoDB!")
}

func newMongoDBConnection(dbConfig *common.DbConfig) (*mongo.Client, context.Context, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(dbConfig.URI).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return client, ctx, nil
}
