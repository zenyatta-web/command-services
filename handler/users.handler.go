package handler

import (
	"context"
	"log"
	"zenyatta-web/command-services/data/models"
	"zenyatta-web/command-services/data/repository"
	proto "zenyatta-web/command-services/proto"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (u *UsersServiceServer) HolaMundo(ctx context.Context, req *proto.HolaMundoRequest) (*proto.HolaMundoResponse, error) {
	return &proto.HolaMundoResponse{Message: "Hola mundo"}, nil
}

func (u *UsersServiceServer) CreateUser(ctx context.Context, request *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	// Imprime el nombre del usuario recibido en el log
	log.Printf("Usuario recibido: %v", request.Name)

	// Crea una instancia de UserModel a partir de los datos de la solicitud
	user := &models.UserModel{
		IdAuth0: request.IdAuth0,
		Name:    request.Name,
		Status:  request.Status,
	}

	// Guarda el usuario en la base de datos.
	user, err := u.userDatabase.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// Prepara la respuesta
	response := &proto.CreateUserResponse{
		Id:      user.Id.Hex(),
		IdAuth0: user.IdAuth0,
		Name:    user.Name,
		Status:  user.Status,
	}

	return response, nil
}

func (u *UsersServiceServer) UpdateUser(ctx context.Context, request *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	// Imprime el nombre del usuario recibido en el log
	log.Printf("Usuario recibido: %v, id: %v", request.Name, request.Id)

	// Verifica si el ID está vacío y lanza un error gRPC de tipo NotFound
	if request.Id == "" {
		return nil, status.Errorf(codes.NotFound, "User ID not provided")
	}

	// Convertir el id de string a ObjectID
	userId, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		log.Printf("error en el handler")
		return nil, grpc.Errorf(codes.InvalidArgument, "Invalid ID format")
	}

	// Crea una instancia de UserModel a partir de los datos de la solicitud
	user := &models.UserModel{
		Id:      userId,
		IdAuth0: request.IdAuth0,
		Name:    request.Name,
		Status:  request.Status,
	}

	// Actualiza el usuario en la base de datos.
	user, err = u.userDatabase.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// Prepara la respuesta
	response := &proto.UpdateUserResponse{
		Id:      user.Id.Hex(),
		IdAuth0: user.IdAuth0,
		Name:    user.Name,
		Status:  user.Status,
	}

	return response, nil
}
