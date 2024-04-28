package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// model
type Note struct {
	/*no space between _id and omitempty else _id generated would be zeros
	and won't allow another insertion, so as to prompt go drivers for
	random _id to be generated*/
	Id          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title,omitempty"`
	Description string             `json:"description" bson:"description,omiempty"`
}

var collection *mongo.Collection

const uri = "mongodb://localhost:27017"

// handler
func tautPadHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	switch request.Method {
	case "POST":
		var note Note

		_ = json.NewDecoder(request.Body).Decode(&note)

		output, err := collection.InsertOne(context.Background(), note)

		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(response).Encode(output)

	case "GET":
		var notes []Note

		output, err := collection.Find(context.TODO(), bson.D{})

		if err != nil {
			panic(err)
		}

		err = output.All(context.TODO(), &notes)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(response).Encode(notes)

	}
}

func getUpdateDeleteTautPadHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	var note Note

	ctx := context.TODO()

	switch request.Method {
	case "GET":
		getId := strings.TrimPrefix(request.URL.Path, "/notes/")

		id, err := primitive.ObjectIDFromHex(getId)

		if err != nil {
			panic(err)
		}

		err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&note)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(response).Encode(note)

	case "PUT":
		err := json.NewDecoder(request.Body).Decode(&note)
		if err != nil {
			panic(err)
		}
		getId := strings.TrimPrefix(request.URL.Path, "/notes/")
		id, err := primitive.ObjectIDFromHex(getId)
		if err != nil {
			panic(err)
		}
		result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": note})
		if err != nil {
			panic(err)
		}
		json.NewEncoder(response).Encode(result)

	case "DELETE":
		getID := strings.TrimPrefix(request.URL.Path, "/notes/")
		id, err := primitive.ObjectIDFromHex(getID)
		if err != nil {
			panic(err)
		}

		output, err := collection.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil {
			panic(err)
		}
		json.NewEncoder(response).Encode(output)
	}
}

// main.go
func main() {

	// client; db connect
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	// collection
	collection = client.Database("tautpad").Collection("notes")

	fmt.Println("Database connected and pinged")

	// mux
	mux := http.NewServeMux()
	mux.Handle("/notes", http.HandlerFunc(tautPadHandler))
	mux.Handle("/notes/", http.HandlerFunc(getUpdateDeleteTautPadHandler))

	// http server
	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
