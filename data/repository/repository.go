package repository

import (
	"context"
	"zenyatta-web/command-services/data/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.UserModel) (*models.UserModel, error)
	UpdateUser(ctx context.Context, user *models.UserModel) (*models.UserModel, error)
}
