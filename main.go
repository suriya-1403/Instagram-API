package main

import (
	"Instagram-API/database_Connect"
	"Instagram-API/models"
	"Instagram-API/routing"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"log"
	"net/http"
	"time"
)

// Creating hash for ciper conversion
func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// encrypt function to encrypt the password
func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

//decrypt function to decrypt the password
func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

// storing the user, post collection from the database for manipulations
var collectionUser, collectionPost = database_Connect.ConnectDB()

// getUsers function is to view all the users of the user collection which is GET method
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

// getUser is to display particular user where id is given in the url which is GET method
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

// createUser is to create user and also in this func password is also encrypted which is POST methos
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
	filter := bson.D{{"email", user.Email}}
	ciper := encrypt([]byte(user.Password), "password")
	newPassEncry := bson.D{
		{"$set", bson.D{
			{"password", ciper},
		}},
	}
	res, err := collectionUser.UpdateOne(context.TODO(), filter, newPassEncry)
	if err != nil {
		database_Connect.GetError(err, w)
	}
	updatedObject := *res
	fmt.Printf("The matched count is : %d, the modified count is : %d", updatedObject.MatchedCount, updatedObject.ModifiedCount)
	json.NewEncoder(w).Encode(result)
}

// createPost is to create the post and also store the timestamp of the post which is POST method
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

// getPost is to display the post from the post id which is a GET method
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

// getPosts is to display all the post of a particular user which is GET method
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

/*
Here I haven't installed any package from the internet except:
	1. monogoDB go driver
	2. standard library of golang
but for routing I have implemented the router in routing folder just like Mux package
where all are standard libraries.
*/
func main() {
	mux := routing.NewRouter() // creating the Router
	mux.HandleFunc("/users", getUsers).Methods("GET")
	mux.HandleFunc("/users/{id}", getUser).Methods("GET")
	mux.HandleFunc("/users", createUser).Methods("POST")
	mux.HandleFunc("/posts", CreatePost).Methods("POST")
	mux.HandleFunc("/posts/{id}", getPost).Methods("GET")
	mux.HandleFunc("/posts/users/{uid}", getPosts).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
