package jokes

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	mux "github.com/gorilla/mux"
)

//Router struct
type Router struct {
}

//ConfigRouter for jokes
func (r *Router) ConfigRouter(router *mux.Router) {

	router.HandleFunc("/", getJokes).Methods("GET")

}

func getJokes(w http.ResponseWriter, r *http.Request) {

	//fmt.Fprintf(w, "This is a subrouter for jokes")

	// response, err := http.Get("http://pokeapi.co/api/v2/pokedex/kanto/")

	response, err := http.Get("https://icanhazdadjoke.com/")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, string(responseData))

}
