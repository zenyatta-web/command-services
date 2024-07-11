package handler

import (
	"context"
	"log"
	"zenyatta-web/command-services/data/models"
	"zenyatta-web/command-services/data/repository"
	proto "zenyatta-web/command-services/proto"
)

// Definicion del servicio, construido en el archivo .proto.
type UsersServiceServer struct {
	proto.UnimplementedUsersServiceServer
	userDatabase repository.UserRepositoryDatabase
}

func ConstructorUsersServiceServer(userDatabase repository.UserRepositoryDatabase) *UsersServiceServer {
	log.Printf("userDatabase recibido: %v", userDatabase)

	return &UsersServiceServer{
		userDatabase: userDatabase,
	}
}

// Implementar la funcion hola mundo.
func (u *UsersServiceServer) HolaMundo(ctx context.Context, req *proto.HolaMundoRequest) (*proto.HolaMundoResponse, error) {
	return &proto.HolaMundoResponse{Message: "Hola mundo"}, nil
}

func (u *UsersServiceServer) CreateUser(ctx context.Context, request *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	// Imprime el nombre del usuario recibido en el log
	log.Printf("Usuario recibido: %v", request.Name)

	// Crea una instancia de UserModel a partir de los datos de la solicitud
	user := &models.UserModel{
		Id:      "",
		IdAuth0: request.IdAuth0,
		Name:    request.Name,
		Status:  request.Status,
	}

	// Guarda el usuario en el repositorio
	user, err := u.userDatabase.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// Prepara la respuesta
	response := &proto.CreateUserResponse{
		Id:      user.Id, // Asegúrate de tener un campo Id en UserModel
		IdAuth0: user.IdAuth0,
		Name:    user.Name,
		Status:  user.Status,
	}

	return response, nil
}
