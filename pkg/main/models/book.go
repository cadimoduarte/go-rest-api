package models

import (
	"context"
	"log"

	"github.com/cadimoduarte/go-rest-api/pkg/main/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//Create Struct
type Book struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Isbn   string             `json:"isbn,omitempty" bson:"isbn,omitempty"`
	Title  string             `json:"title" bson:"title,omitempty"`
	Author *Author            `json:"author" bson:"author,omitempty"`
}

type Author struct {
	FirstName string `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

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
