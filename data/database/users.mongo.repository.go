package database

import (
	"context"
	"zenyatta-web/command-services/data/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userMongoRepository struct {
	usersCollection *mongo.Collection
}

func ConstructorUserMongoRepository(usersCollection *mongo.Collection) *userMongoRepository {
	return &userMongoRepository{usersCollection}
}

func (r *userMongoRepository) CreateUser(ctx context.Context, user *models.UserModel) (*models.UserModel, error) {
	result, err := r.usersCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	userId := result.InsertedID
	err = r.usersCollection.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userMongoRepository) UpdateUser(ctx context.Context, user *models.UserModel) (*models.UserModel, error) {
	filter := bson.M{"_id": user.Id}
	update := bson.M{"$set": user}

	_, err := r.usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	err = r.usersCollection.FindOne(ctx, bson.M{"_id": user.Id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}