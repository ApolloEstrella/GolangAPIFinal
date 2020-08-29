package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, //you service is available and allowed for this base url
		AllowedMethods: []string{
			http.MethodGet, //http methods for your app
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},

		AllowedHeaders: []string{
			"*", //or you can your header key values which you are using in your application

		},
	})

	// create handler
	handler := c.Handler(router)

	router.HandleFunc("/api/person", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/api/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/api/person/{id}", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/api/delete/{id}", DeletePersonEndpoint).Methods("DELETE")
	router.HandleFunc("/api/update/{id}", UpdatePersonEndpoint).Methods("PUT")
	http.ListenAndServe(":12345", handler)
}
