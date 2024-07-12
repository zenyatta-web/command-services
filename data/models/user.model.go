package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserModel struct {
	Id      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	IdAuth0 string             `bson:"id_auth0" json:"id_auth0"`
	Name    string             `bson:"name" json:"name"`
	Status  bool               `bson:"status" json:"status"`
}
