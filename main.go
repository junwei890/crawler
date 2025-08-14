package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/junwei890/se-cli/src"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(errors.New("couldn't load .env file"))
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal(errors.New("mondodb uri not set"))
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database("se-cli")
	collection := db.Collection("content")

	if err := src.Init(collection); err != nil {
		log.Fatal(err)
	}
}
