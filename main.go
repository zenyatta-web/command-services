package main

import (
	"log"
	"net"
	"zenyatta-web/command-services/config"
	"zenyatta-web/command-services/data/database"
	"zenyatta-web/command-services/handler"

	pb "zenyatta-web/command-services/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Comienza cargando la configuración necesaria del servicio.
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Conexión a MongoDB
	mongoConfig := config.Mongo()
	db, err := database.NewDatabase(mongoConfig.URI, mongoConfig.Database)
	if err != nil {
		log.Fatalf("Error conectando a MongoDB: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Error closing MongoDB connection: %v", err)
		}
	}()

	proveedor := database.ConstructorUserMongoRepository(db.Database.Collection("users"))
	log.Printf("Proveedor: %v", proveedor)
	userServiceServer := handler.ConstructorUsersServiceServer(proveedor)

	// Inicialización del servidor gRPC
	lis, err := net.Listen("tcp", config.Address())
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %v", config.Address(), err)
	}

	server := grpc.NewServer()
	pb.RegisterUsersServiceServer(server, userServiceServer)

	// Habilitar la reflexión
	reflection.Register(server)

	log.Printf("Starting gRPC server on %v", config.Address())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
