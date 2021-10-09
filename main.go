package main

import (
	"Instagram-API/database_Connect"
	"Instagram-API/models"
	"Instagram-API/routing"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

var collectionUser, collectionPost = database_Connect.ConnectDB()

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []models.User
	cur, err := collectionUser.Find(context.TODO(), bson.M{})

	if err != nil {
		database_Connect.GetError(err, w)
		return
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var user models.User
		err := cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(users)
}
func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	var params = routing.Vars(r)
	ID, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": ID}
	err := collectionUser.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		database_Connect.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	result, err := collectionUser.InsertOne(context.TODO(), user)

	if err != nil {
		database_Connect.GetError(err, w)
		return
	}
	fmt.Println("Created!!!")
	json.NewEncoder(w).Encode(result)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	var post models.Post
	_ = json.NewDecoder(r.Body).Decode(&post)
	result, err := collectionPost.InsertOne(context.TODO(), post)
	if err != nil {
		database_Connect.GetError(err, w)
		return
	}
	fmt.Println("Post Created!!!")
	filter := bson.D{{"_uid", post.ID}}
	newTime := bson.D{
		{"$set", bson.D{
			{"time", primitive.Timestamp{T: uint32(time.Now().Unix())}},
		}},
	}
	res, err := collectionPost.UpdateOne(context.TODO(), filter, newTime)
	if err != nil {
		log.Fatal(err)
	}
	updatedObject := *res
	fmt.Printf("The matched count is : %d, the modified count is : %d", updatedObject.MatchedCount, updatedObject.ModifiedCount)
	json.NewEncoder(w).Encode(result)
}

func getPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var post models.Post
	var params = routing.Vars(r)
	ID, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": ID}
	err := collectionPost.FindOne(context.TODO(), filter).Decode(&post)
	if err != nil {
		database_Connect.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(post)
}
func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = routing.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["uid"])
	filter := bson.M{"_uid": id}

	cur, err := collectionPost.Find(context.TODO(), filter)
	if err != nil {
		database_Connect.GetError(err, w)
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var post models.Post
		err := cur.Decode(&post)
		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(post.ImageURL)
	}

}
func main() {
	mux := routing.NewRouter()
	mux.HandleFunc("/users", getUsers).Methods("GET")
	mux.HandleFunc("/users/{id}", getUser).Methods("GET")
	mux.HandleFunc("/users", createUser).Methods("POST")
	mux.HandleFunc("/posts", CreatePost).Methods("POST")
	mux.HandleFunc("/posts/{id}", getPost).Methods("GET")
	mux.HandleFunc("/posts/users/{uid}", getPosts).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
