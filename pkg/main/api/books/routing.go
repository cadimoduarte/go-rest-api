package books

import (
	"encoding/json"
	"fmt"

	"github.com/cadimoduarte/go-rest-api/pkg/main/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

//FastRouter struct
type FastRouter struct {
	MongoClient *mongo.Client
	Router      *routing.Router
}

var mongoDBClient *mongo.Client

//StartRouter for books
func (r *FastRouter) StartRouter() {

	fmt.Print("Config fast 'book' subroutes...")

	mongoDBClient = r.MongoClient

	books := r.Router.Group("/api/books")

	books.Get("/", listAll)
	books.Get("/<id>", get)
	books.Post("", create)
	books.Put("/<id>", update)
	books.Delete("/<id>", delete)

	fmt.Println(" Done")

}

func listAll(ctx *routing.Context) error {
	var book = models.Book{}

	books, err := book.Load(mongoDBClient)

	if err != nil {
		ctx.SetContentType("application/json")
		ctx.Response.Header.Set("Content-Type", "application/json;charset=UTF-8")
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "Problem loading books")
		return err
	}

	ctx.SetContentType("application/json")
	ctx.Response.Header.Set("Content-Type", "application/json;charset=UTF-8")
	ctx.SetStatusCode(fasthttp.StatusOK)
	json.NewEncoder(ctx).Encode(books)

	return nil
}

func get(ctx *routing.Context) error {
	id, _ := primitive.ObjectIDFromHex(ctx.Param("id"))

	var book = models.Book{
		ID: id,
	}

	err := book.Get(mongoDBClient)

	if err != nil {
		ctx.SetContentType("application/json")
		ctx.Response.Header.Set("Content-Type", "application/json;charset=UTF-8")
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "Problem loading a book")
		return err
	}

	ctx.SetContentType("application/json")
	ctx.Response.Header.Set("Content-Type", "application/json;charset=UTF-8")
	ctx.SetStatusCode(fasthttp.StatusOK)
	json.NewEncoder(ctx).Encode(book)

	return nil

}

func create(ctx *routing.Context) error {

	var book models.Book

	// _ = json.NewDecoder(ctx.PostBody()).Decode(&book)
	_ = json.Unmarshal(ctx.PostBody(), &book)

	err := book.Insert(mongoDBClient)

	if err != nil {
		ctx.SetContentType("application/json")
		ctx.Response.Header.Set("Content-Type", "application/json;charset=UTF-8")
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "Problem creating a book")
		return err
	}

	ctx.SetContentType("application/json")
	ctx.Response.Header.Set("Content-Type", "application/json;charset=UTF-8")
	ctx.SetStatusCode(fasthttp.StatusCreated)
	json.NewEncoder(ctx).Encode(book)

	return nil
}

func update(ctx *routing.Context) error {

	id, _ := primitive.ObjectIDFromHex(ctx.Param("id"))

	book := models.Book{ID: id}

	_ = json.Unmarshal(ctx.PostBody(), &book)

	err := book.Update(mongoDBClient)

	if err != nil {
		ctx.SetContentType("application/json")
		ctx.Response.Header.Set("Content-Type", "application/json;charset=UTF-8")
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "Problem updating a book")
		return err
	}

	ctx.SetContentType("application/json")
	ctx.Response.Header.Set("Content-Type", "application/json;charset=UTF-8")
	ctx.SetStatusCode(fasthttp.StatusOK)
	json.NewEncoder(ctx).Encode(book)

	return nil
}

func delete(ctx *routing.Context) error {

	id, _ := primitive.ObjectIDFromHex(ctx.Param("id"))

	book := models.Book{ID: id}

	deleteErr := book.Delete(mongoDBClient)

	if deleteErr != nil {
		ctx.SetContentType("application/json")
		ctx.Response.Header.Set("Content-Type", "application/json;charset=UTF-8")
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "Problem loading a book")
		return deleteErr
	}

	ctx.SetContentType("application/json")
	ctx.Response.Header.Set("Content-Type", "application/json;charset=UTF-8")
	ctx.SetStatusCode(fasthttp.StatusOK)

	return nil
}
