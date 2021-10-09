package main

import (
	"Instagram-API/database_Connect"
	"Instagram-API/models"
	"Instagram-API/routing"
	"context"
	"encoding/json"
	"fmt"
	bson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

var collectionUser, collectionPost = database_Connect.ConnectDB()

type Post_main struct {
	IdMain              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UidMain             primitive.ObjectID `json:"_uid,omitempty" bson:"_uid,omitempty"`
	CaptionMain         string             `json:"caption,omitempty" bson:"caption,omitempty"`
	ImageurlMain        *imageurlMain      `json:"imageURL,omitempty" bson:"imageURL,omitempty"`
	PostedtimestampMain time.Time          `json:"time,omitempty" bson:"time"`
}
type imageurlMain struct {
	URl string `json:"url,omitempty" bson:"url,omitempty"`
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []models.User
	cur, err := collectionUser.Find(context.TODO(), bson.M{})

	if err != nil {
		database_Connect.GetError(err, w)
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

func CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	var post models.Post
	// we decode our body request params
	_ = json.NewDecoder(r.Body).Decode(&post)
	//post = bson.D{
	//		{"time",  primitive.Timestamp{T:uint32(time.Now().Unix())}},
	//	}
	// insert our book model.
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
	//updateTime(w,r)
	//result = context.TODO()
	json.NewEncoder(w).Encode(result)

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
func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = routing.Vars(r)
	//var posts []models.Post
	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(params["uid"])

	//var post models.Post

	// Create filter
	filter := bson.M{"_uid": id}

	cur, err := collectionPost.Find(context.TODO(), filter)
	if err != nil {
		database_Connect.GetError(err, w)
		return
	}

	// Close the cursor once finished
	/*A defer statement defers the execution of a function until the surrounding function returns.
	simply, run cur.Close() process but after cur.Next() finished.*/
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var post models.Post
		// & character returns the memory address of the following variable.
		err := cur.Decode(&post) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}
		//elem := &Post_main{}

		// add item our array
		//posts = append(posts, post)
		json.NewEncoder(w).Encode(post.ImageURL)
	}

	//if err := cur.Err(); err != nil {
	//	log.Fatal(err)
	//}

	//json.NewEncoder(w).Encode(posts)
	//
	//elem := &Post_main{}
	//err := resultDoc.Decode(elem)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//json.NewEncoder(w).Encode(elem.ImageurlMain)

}
func main() {
	mux := routing.NewRouter()
	mux.HandleFunc("/api/users", getUsers).Methods("GET")
	mux.HandleFunc("/api/users/{id}", getUser).Methods("GET")
	mux.HandleFunc("/api/users", createUser).Methods("POST")
	mux.HandleFunc("/api/posts", CreatePost).Methods("POST")
	mux.HandleFunc("/api/posts/{id}", getPost).Methods("GET")
	mux.HandleFunc("/api/posts/users/{uid}", getPosts).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
