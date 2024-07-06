package models

type UserModel struct {
	Id      string `bson:"_id,omitempty" json:"id"`
	IdAuth0 string `bson:"id_auth0" json:"id_auth0"`
	Name    string `bson:"name" json:"name"`
	Status  bool   `bson:"status" json:"status"`
}
