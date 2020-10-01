package books

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cadimoduarte/go-rest-api/pkg/main/helper"
	"github.com/cadimoduarte/go-rest-api/pkg/main/models"
	mux "github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//Router struct
type Router struct {
	MongoClient *mongo.Client
}

var mongoClient *mongo.Client

//ConfigRouter for books
func (r *Router) ConfigRouter(router *mux.Router) {

	mongoClient = r.MongoClient

	router.HandleFunc("/", getBooks).Methods("GET")
	router.HandleFunc("/{id}", getBook).Methods("GET")
	router.HandleFunc("", createBook).Methods("POST")
	router.HandleFunc("/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/{id}", deleteBook).Methods("DELETE")

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
