package main

/*
	Este archivo define una función llamada newTracerProvider que retorna un proveedor de
	trazas de OpenTelemetry configurado para utilizar el exportador Jaeger,
	el cual enviará los registros de trazas (spans) a la URL proporcionada.
*/
import (
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func newTracerProvider(serviceName, serviceID, url string) (*tracesdk.TracerProvider, error) {
	/*Esta función newTracerProvider crea y configura un exportador Jaeger
	 	para enviar trazas a la URL especificada. Luego, utiliza este exportador para
		inicializar un proveedor de trazas TracerProvider. Además, agrega información
		sobre el servicio como atributos al recurso asociado con el proveedor de
		trazas, como el nombre del servicio, la versión del servicio y
		el ID de instancia del servicio.*/
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(version),
			semconv.ServiceInstanceIDKey.String(serviceID),
		)),
	)
	return tp, nil
}
