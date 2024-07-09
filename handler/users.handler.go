package handler

import (
	"context"
	proto "zenyatta-web/command-services/proto"
)

// Definicion del servicio, construido en el archivo .proto.
type UsersServiceServer struct {
	proto.UnimplementedUsersServiceServer
}

// Implementar la funcion hola mundo.
func (u *UsersServiceServer) HolaMundo(ctx context.Context, req *proto.HolaMundoRequest) (*proto.HolaMundoResponse, error) {
	return &proto.HolaMundoResponse{Message: "Hola mundo"}, nil
}
