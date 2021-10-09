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
			database_Connect.GetError(err,w)
			return
		}

		// Close the cursor once finished
		/*A defer statement defers the execution of a function until the surrounding function returns.
		simply, run cur.Close() process but after cur.Next() finished.*/
		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {

			// create a value into which the single document can be decoded
			var user models.User
			// & character returns the memory address of the following variable.
			err := cur.Decode(&user) // decode similar to deserialize process.
			if err != nil {
				log.Fatal(err)
			}

			// add item our array
			users = append(users, user)
		}

		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(users) // encode similar to serialize process.
	//} else {
	//	fmt.Print("Not Get")
	//}
}

//func ObjectIDFromHex(s string) (ObjectID, error)
func getUser(w http.ResponseWriter, r *http.Request) {
	// set header.
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	//id := r.URL.Path[len("/user/"):]
	//id := Vars(r)
	var params = routing.Vars(r)

	// string to primitive.ObjectID
	ID, _ := primitive.ObjectIDFromHex(params["id"])
	//ID, _ := primitive.ObjectIDFromHex(id)
	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
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

	// we decode our body request params
	_ = json.NewDecoder(r.Body).Decode(&user)

	// insert our book model.
	result, err := collectionUser.InsertOne(context.TODO(), user)

	if err != nil {
		database_Connect.GetError(err, w)
		return
	}
	fmt.Println("Created!!!")
	json.NewEncoder(w).Encode(result)
}

func CreatePost(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Context-Type","application/json")
	var post models.Post
	// we decode our body request params
	_ = json.NewDecoder(r.Body).Decode(&post)

	// insert our book model.
	result, err := collectionPost.InsertOne(context.TODO(), post)

	if err != nil {
		database_Connect.GetError(err, w)
		return
	}
	fmt.Println("Post Created!!!")
	filter := bson.M{"_id": post.ID}
	ti := time.Now()
	_ = json.NewDecoder(r.Body).Decode(&post)
	json.NewEncoder(w).Encode(result)
	// prepare update model.
	update := bson.D{
		{"$set", bson.D{
			{"_uid", post.UID},
			{"caption", post.Caption},
			{"imageURL", bson.D{
				{"url", post.ImageURL.URl},
			}},
			{"time",ti},
		}},
	}

	err = collectionPost.FindOneAndUpdate(context.TODO(), filter, update).Decode(&post)

	if err != nil {
		database_Connect.GetError(err, w)
		return
	}

	//post.ID =

	//json.NewEncoder(w).Encode(post)
}

func getPost(w http.ResponseWriter, r *http.Request) {
	// set header.
	w.Header().Set("Content-Type", "application/json")

	var post models.Post

	var params = routing.Vars(r)

	// string to primitive.ObjectID
	ID, _ := primitive.ObjectIDFromHex(params["id"])
	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	filter := bson.M{"_id": ID}
	err := collectionPost.FindOne(context.TODO(), filter).Decode(&post)

	if err != nil {
		database_Connect.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(post)
}
func getPosts(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")

	var params = routing.Vars(r)

	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(params["uid"])

	var post models.Post

	// Create filter
	filter := bson.M{"_uid": id}

	//ti := time.Now()
	//_ = json.NewDecoder(r.Body).Decode(&post)
	//json.NewEncoder(w).Encode(result)

	err := collectionPost.FindOne(context.TODO(), filter).Decode(&post)

	if err != nil {
		database_Connect.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(post)
	//if err != nil {
	//	database_Connect.GetError(err, w)
	//	return
	//}
	//
	//// Close the cursor once finished
	//defer cur.Close(context.TODO())
	//
	//for cur.Next(context.TODO()) {
	//
	//	// create a value into which the single document can be decoded
	//	var post models.Post
	//	// & character returns the memory address of the following variable.
	//	err := cur.Decode(&post) // decode similar to deserialize process.
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	// add item our array
	//	posts = append(posts, post)
	//}

	//if err := cur.Err(); err != nil {
	//	log.Fatal(err)
	//}
	//
	//json.NewEncoder(w).Encode(post) // encode similar to serialize process.
}
func main(){
	mux := routing.NewRouter()
	mux.HandleFunc("/api/users", getUsers).Methods("GET")
	mux.HandleFunc("/api/users/{id}", getUser).Methods("GET")
	mux.HandleFunc("/api/users", createUser).Methods("POST")
	mux.HandleFunc("/api/posts",CreatePost).Methods("POST")
	mux.HandleFunc("/api/posts/{id}", getPost).Methods("GET")
	mux.HandleFunc("/api/posts/users/{uid}", getPosts).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", mux))
}


