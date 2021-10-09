package database_Connect

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

//clientOptions := options.Client().
//
//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//defer cancel()
//client, err := mongo.Connect(ctx, clientOptions)
//if err != nil {
//log.Fatal(err)
//}

// ConnectDB : This is helper function to connect mongoDB
// If you want to export your function. You must to start upper case function name. Otherwise you won't see your function when you import that on other class.
func ConnectDB() (*mongo.Collection, *mongo.Collection) {

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb+srv://root:root@appointy.scljs.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection_user := client.Database("Instagram_API").Collection("Users")
	collection_post := client.Database("Instagram_API").Collection("Post")


	return collection_user, collection_post
}

// ErrorResponse : This is error model.
type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

// GetError : This is helper function to prepare error model.
// If you want to export your function. You must to start upper case function name. Otherwise you won't see your function when you import that on other class.
func GetError(err error, w http.ResponseWriter) {

	log.Fatal(err.Error())
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}

// Configuration model
type Configuration struct {
	Port             string
	ConnectionString string
}
