package handler

/*
	Este paquete se encarga de manejar las verificaciones de salud (Health Check) de un microservicio.
	Los Health Checks son esenciales para monitorizar el estado de los servicios y garantizar que
	estén funcionando correctamente. Este código define dos métodos, Check y Watch, que responden a
	solicitudes de verificación de salud.
*/

import (
	"context" //Para manejar el contexto de las solicitudes gRPC.

	pb "github.com/go-micro/demo/productcatalogservice/proto"
	"google.golang.org/grpc/codes"  //Para manejar los códigos de estado y errores en gRPC.
	"google.golang.org/grpc/status" //Para manejar los códigos de estado y errores en gRPC.
)

type Health struct{} //La estructura Health actúa como el receptor para los métodos de verificación de salud.

// Este método simplemente asigna el estado SERVING a la respuesta, indicando que el servicio está operativo.
func (h *Health) Check(ctx context.Context, req *pb.HealthCheckRequest, rsp *pb.HealthCehckResponse) error {
	rsp.Status = pb.HealCheckResponse_SERVING
	return nil
}

/*
Este método devuelve un error con el código Unimplemented, indicando que la funcionalidad de verificación
de salud vía Watch no está implementada.
*/
func (h *Health) Watch(ctx context.Context, req *pb.HealthCheckRequest, stream pb.Health_WatchStream) error {
	return status.Errorf(codes.Unimplemented, "health check via Watch not implemented")
}
