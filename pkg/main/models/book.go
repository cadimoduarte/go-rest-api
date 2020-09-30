package models

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/cadimoduarte/go-rest-api/pkg/main/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
func (b *Book) Load(client *mongo.Client) (books []Book, err error) {
	db := client.Database(config.MongoDbDatabase())
	collection := db.Collection("books")

	ctx := context.Background()

	cur, err := collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cur.Close(context.TODO())

	var book Book

	for cur.Next(context.TODO()) {

		err := cur.Decode(&book)
		if err != nil {
			log.Fatal(err)
		}
		books = append(books, book)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
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

	fmt.Println("InsertOne() result type: ", reflect.TypeOf(result))
	fmt.Println("InsertOne() API result:", result)
	newID := result.InsertedID
	fmt.Println("InsertOne() newID:", newID)
	fmt.Println("InsertOne() newID type:", reflect.TypeOf(newID))

	b.ID = newID.(primitive.ObjectID)

	return nil
}
