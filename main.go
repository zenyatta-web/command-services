package main

import (
	"context"
	"strings"
	"time"
	"zenyatta-web/command-services/config"
	"zenyatta-web/command-services/data/database"
	"zenyatta-web/command-services/handler"

	pb "zenyatta-web/command-services/proto"

	grpcc "github.com/go-micro/plugins/v4/client/grpc"
	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
	grpcs "github.com/go-micro/plugins/v4/server/grpc"
	"github.com/go-micro/plugins/v4/wrapper/trace/opentelemetry"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/server"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// Se definen variables para el nombre y la versión del servicio.
var (
	name    = "productcatalogservice"
	version = "1.0.0"
)

func main() {
	// Comienza cargando la configuración necesaria del servicio.
	if err := config.Load(); err != nil {
		logger.Fatal(err)
	}

	// Conexión a MongoDB
	mongoConfig := config.Mongo()
	db, err := database.NewDatabase(mongoConfig.URI, mongoConfig.Database)
	if err != nil {
		logger.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Fatal("Error closing MongoDB connection: %v", err)
		}
	}()

	// Crear una instancia del servicio utilizando micro.NewService()
	srv := micro.NewService(
		micro.Server(grpcs.NewServer()),
		micro.Client(grpcc.NewClient()),
	)
	opts := []micro.Option{
		micro.Name(name),
		micro.Version(version),
		micro.Address(config.Address()),
	}

	// Configuración de trazabilidad
	if cfg := config.Tracing(); cfg.Enable {
		tp, err := newTracerProvider(name, srv.Server().Options().Id, cfg.Jaeger.URL)
		if err != nil {
			logger.Fatal(err)
		}
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			if err := tp.Shutdown(ctx); err != nil {
				logger.Fatal(err)
			}
		}()
		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
		traceOpts := []opentelemetry.Option{
			opentelemetry.WithHandleFilter(func(ctx context.Context, r server.Request) bool {
				if e := r.Endpoint(); strings.HasPrefix(e, "Health.") {
					return true
				}
				return false
			}),
		}
		opts = append(opts, micro.WrapHandler(opentelemetry.NewHandlerWrapper(traceOpts...)))
	}

	// Inicialización del servicio con las opciones configuradas.
	srv.Init(opts...)

	// Registro de los controladores de los servicios de catálogo de productos y de salud
	if err := pb.RegisterProductCatalogServiceHandler(srv.Server(), new(handler.ProductCatalogService)); err != nil { // Pasar la conexión a MongoDB al handler
		logger.Fatal(err)
	}
	if err := pb.RegisterHealthHandler(srv.Server(), new(handler.Health)); err != nil {
		logger.Fatal(err)
	}

	// Ejecutar el servicio
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
