package mongodb

import (
	"context"

	"github.com/thienkb1123/go-clean-arch/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri string = "mongodb://localhost:27017"

func New(c *config.MySQLConfig) (*mongo.Client, error)  {
    serverAPI := options.ServerAPI(options.ServerAPIVersion1)
    opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

    client, err := mongo.Connect(context.TODO(), opts)
    if err != nil {
        panic(err)
    }
    err = client.Ping(context.TODO(), nil)
    mongoClient := client

    return mongoClient, err
}

