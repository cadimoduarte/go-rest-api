package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cadimoduarte/go-rest-api/pkg/main/config"
	"github.com/cadimoduarte/go-rest-api/pkg/main/helper"
	"github.com/cadimoduarte/go-rest-api/pkg/main/models"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome! It`s my first Go Web API application!")
}

var mongoTries int

var mongoClient *mongo.Client

func init() {

	mongoTries = 0

	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found")
	}
	os.Setenv("TZ", "America/Sao_Paulo")
	os.Setenv("LOG_APP", "MAIN")

	mongoClient = connectOnMongo()

}

func connectOnMongo() *mongo.Client {
	log.Println("Connecting on MongoDB...")
	opt := &options.ClientOptions{}
	opt = opt.ApplyURI(config.MongoDbURI())
	opt = opt.SetMaxPoolSize(config.MongoDbMaxPoolSize())
	client, err := mongo.NewClient(opt)
	if err != nil {
		log.Panic(err)
	}

	err = client.Connect(context.TODO())

	if err != nil {
		time.Sleep(time.Second * 1)
		log.Println(err.Error())
		mongoTries++
		if mongoTries > config.MongoDbMaxRetries() {
			log.Panic(err)
		}
		log.Println("Trying again...")
		connectOnMongo()
	} //else {
	//app.Resources["mongo"] = client
	//}
	return client
}

func main() {

	//Init Router
	router := mux.NewRouter()
	router.HandleFunc("/", homeLink)

	// arrange our route
	router.HandleFunc("/api/books", getBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/books", createBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	// set our port address
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book = models.Book{}

	books, err := book.Load(mongoClient)

	if err != nil {
		//TODO: create new error handlers
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(books) // encode similar to serialize process.
}

func getBook(w http.ResponseWriter, r *http.Request) {
	// set header.
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var book = models.Book{
		ID: id,
	}

	err := book.Get(mongoClient)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(book)

}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book models.Book

	_ = json.NewDecoder(r.Body).Decode(&book)

	err := book.Insert(mongoClient)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(params["id"])

	book := models.Book{ID: id}

	_ = json.NewDecoder(r.Body).Decode(&book)

	err := book.Update(mongoClient)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	fmt.Println(book)

	json.NewEncoder(w).Encode(book)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	// Set header
	w.Header().Set("Content-Type", "application/json")

	// get params
	var params = mux.Vars(r)

	// string to primitve.ObjectID
	id, _ := primitive.ObjectIDFromHex(params["id"])

	book := models.Book{ID: id}

	deleteErr := book.Delete(mongoClient)

	if deleteErr != nil {
		return
	}
}
