package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Person ...
type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

// CreatePersonEndpoint ...
func CreatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var person Person
	_ = json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database("mymongodb01").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}

// GetPeopleEndpoint ...
func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var people []Person
	collection := client.Database("mymongodb01").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(people)
}

// GetPersonEndpoint ...
func GetPersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person Person
	collection := client.Database("mymongodb01").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Person{ID: id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(person)
}

// DeletePersonEndpoint ...
func DeletePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person Person
	collection := client.Database("mymongodb01").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOneAndDelete(ctx, Person{ID: id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(person)
}

// UpdatePersonEndpoint ...
func UpdatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	bodyBytes, _ := ioutil.ReadAll(request.Body)
	body := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	fmt.Println(bodyBytes)
	fmt.Println(body)

	type MongoFields struct {
		FirstName string
		LastName  string
	}

	var bird MongoFields
	json.Unmarshal([]byte(bodyBytes), &bird)
	fmt.Printf("Species: %s, Description: %s", bird.FirstName, bird.LastName)

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person Person
	collection := client.Database("mymongodb01").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	filter := Person{ID: id}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "firstname", Value: bird.FirstName}, {Key: "lastname", Value: bird.LastName},
	}}}

	err := collection.FindOneAndUpdate(ctx, filter, update).Decode(&person)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
	json.NewEncoder(response).Encode(person)
}
