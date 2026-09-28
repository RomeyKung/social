package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// import ()

var Client *mongo.Client
var DB *mongo.Database

func Connect() {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	mongoUri := "mongodb://localhost:27017"
	Client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))

	if err != nil {
		panic(err)
	}

	DB = Client.Database("social")
}
