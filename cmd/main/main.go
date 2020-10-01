package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cadimoduarte/go-rest-api/pkg/main/api/books"
	"github.com/cadimoduarte/go-rest-api/pkg/main/api/jokes"
	"github.com/cadimoduarte/go-rest-api/pkg/main/config"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

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
	}

	return client
}

func main() {

	//Init Router
	router := mux.NewRouter()

	jokesRouter := jokes.Router{}
	jokesRouter.ConfigRouter(router.PathPrefix("/api/jokes").Subrouter())

	booksRouter := books.Router{
		MongoClient: mongoClient,
	}
	booksRouter.ConfigRouter(router.PathPrefix("/api/books").Subrouter())

	router.HandleFunc("/", homeLink)

	// set our port address
	fmt.Println("Starting server")
	log.Fatal(http.ListenAndServe(":8080", router))
}
