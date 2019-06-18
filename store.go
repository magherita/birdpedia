package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// The db variable is a package level variable that will be available for
// use throughout our application code
var db Database

// Database will have two methods, to add a new bird,
// and to get all existing birds
// Each method returns an error, in case something goes wrong
type Database interface {
	CreateBird(bird *Bird) error
	GetBirds() ([]*Bird, error)
}

// Connect struct will implement the `Store` interface
type Connect struct {
	client *mongo.Client
}

// CreateBird ...
func (connect *Connect) CreateBird(bird *Bird) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := connect.client.Database("test").Collection("birds")
	_, err := collection.InsertOne(ctx, bird)
	if err != nil {
		panic(err.Error()) // to do: use a logger
	}
	return err
}

// GetBirds ...
func (connect *Connect) GetBirds() ([]*Bird, error) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	collection := connect.client.Database("test").Collection("birds")
	rows, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close(ctx)

	birds := []*Bird{}
	for rows.Next(ctx) {
		bird := &Bird{}
		err := rows.Decode(&bird)
		if err != nil {
			log.Fatal(err)
		}
		birds = append(birds, bird)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return birds, nil
}

// InitDatabase method will be called to initialize the store
func InitDatabase(database Database) {
	db = database
}
