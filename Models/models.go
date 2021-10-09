package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email" bson:"email,omitempty"`
	Password string             `json:"password" bson:"password,omitempty"`
}

type Posts struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Caption   string             `json:"caption,omitempty" bson:"caption,omitempty"`
	ImageURL  string             `json:"imageurl,omitempty" bson:"imageurl,omitempty"`
	UserID    primitive.ObjectID `json:"userid,omitempty" bson:"userid,omitempty"`
	Timestamp time.Time          `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}

type Page struct {
	Key int `json:"key,omitempty" bson:"key,omitempty"`
}
