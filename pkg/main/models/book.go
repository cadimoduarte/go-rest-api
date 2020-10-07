package models

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/cadimoduarte/go-rest-api/pkg/main/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Book Struct
type Book struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Isbn   string             `json:"isbn,omitempty" bson:"isbn,omitempty"`
	Title  string             `json:"title" bson:"title,omitempty"`
	Author *Author            `json:"author" bson:"author,omitempty"`
}

//Author Struct
type Author struct {
	FirstName string `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

// Load all books from DB
func (b *Book) Load(client *mongo.Client) ([]Book, error) {
	db := client.Database(config.MongoDbDatabase())
	collection := db.Collection("books")

	ctx := context.Background()

	cur, err := collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cur.Close(context.TODO())

	var books []Book
	var book Book

	for cur.Next(context.TODO()) {

		err := cur.Decode(&book)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		books = append(books, book)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return books, nil

}

// Get a book from DB
func (b *Book) Get(client *mongo.Client) error {
	db := client.Database(config.MongoDbDatabase())
	collection := db.Collection("books")

	ctx := context.Background()

	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	filter := bson.M{"_id": b.ID}
	err := collection.FindOne(ctx, filter).Decode(&b)

	if err != nil {
		// helper.GetError(err, w)
		return err
	}

	return nil

}

// Insert a book on DB
func (b *Book) Insert(client *mongo.Client) error {

	db := client.Database(config.MongoDbDatabase())
	collection := db.Collection("books")

	ctx := context.TODO()

	// insert our book model.
	result, err := collection.InsertOne(ctx, b)

	if err != nil {
		fmt.Println("Insert Error:", err)
		return err
	}

	// primitive is an interface, so it needs a type assertion. A type assertion provides access to an interface value's underlying concrete value.
	b.ID = result.InsertedID.(primitive.ObjectID)

	return nil
}

//Update a book on DB
func (b *Book) Update(client *mongo.Client) error {

	db := client.Database(config.MongoDbDatabase())
	collection := db.Collection("books")

	filter := bson.M{"_id": b.ID}

	update := bson.M{
		"$set": bson.M{
			"isbn":  b.Isbn,
			"title": b.Title,
			"author": bson.M{
				"firstname": b.Author.FirstName,
				"lastname":  b.Author.LastName,
			},
		},
	}

	returnDocument := options.After
	updateOptions := options.FindOneAndUpdateOptions{ReturnDocument: &returnDocument}
	err := collection.FindOneAndUpdate(context.TODO(), filter, update, &updateOptions).Decode(&b)

	if err != nil {
		fmt.Println("Update Error:", err)
		return err
	}

	return nil
}

//Delete a book from DB
func (b *Book) Delete(client *mongo.Client) error {

	db := client.Database(config.MongoDbDatabase())
	collection := db.Collection("books")

	filter := bson.M{"_id": b.ID}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		fmt.Println("Delete Error:", err)
		return err
	}

	if deleteResult.DeletedCount == 0 {
		return errors.New("No tasks were deleted")
	}

	return nil
}
