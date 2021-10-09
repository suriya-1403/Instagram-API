package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name,omitempty" bson:"name,omitempty"`
	Email  string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string            `json:"password,omitempty" bson:"password,omitempty"`
}
type Post struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UID    primitive.ObjectID `json:"_uid,omitempty" bson:"_uid,omitempty"`
	Caption string `json:"caption,omitempty" bson:"caption,omitempty"`
	ImageURL   *imageURL `json:"imageURL,omitempty" bson:"imageURL,omitempty"`
	PostedTimeStamp time.Time `json:"time,omitempty" bson:"time"`
}
type imageURL struct {
	URl string `json:"url,omitempty" bson:"url,omitempty"`
}