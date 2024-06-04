# Este Dockerfile está diseñado para crear una imagen Docker para un microservicio escrito en Go.

# 1. Primera Etapa: Builder
# Utiliza la imagen base de Golang 1.19.0 sobre Alpine Linux, lo cual es una versión ligera de Linux. 
# La etiqueta AS builder define un alias para esta etapa de construcción.
FROM golang:1.19.0-alpine AS builder

# Establece variables de entorno para deshabilitar CGO (código C en Go) y especificar 
#   que el sistema operativo objetivo es Linux.
ENV CGO_ENABLED=0 GOOS=linux

# Define el directorio de trabajo dentro del contenedor.
WORKDIR /go/src/hipstershop

# Utiliza el gestor de paquetes apk para instalar certificados SSL, 
#   la herramienta make y protoc (compilador de Protocol Buffers).
RUN apk --update --no-cache add ca-certificates make protoc

# Descarga la herramienta grpc_health_probe utilizada para verificar la salud del servicio gRPC y la hace ejecutable.
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.11 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

# Copia el Makefile y los archivos de dependencias de Go al contenedor.
COPY Makefile go.mod go.sum ./

# Configura el proxy de Go.
RUN go env -w GOPROXY=https://goproxy.io,direct/

# Ejecuta comandos para inicializar el entorno de desarrollo y descargar las dependencias de Go.
RUN make init && go mod download

# Copia todos los archivos del proyecto al contenedor.
COPY . .

# Ejecuta comandos definidos en el Makefile para compilar archivos de Protocol Buffers y limpiar dependencias.
RUN make proto tidy

# Define un argumento para pasar opciones de depuración de skaffold.
ARG SKAFFOLD_GO_GCFLAGS

#Compila el binario de Go con las banderas de compilación para depuración si las hay.
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -o /go/src/hipstershop/productcatalogservice . 

#------------------------------------------------------------------------------------------------------------------------------------------------------#

# 2. Segunda Etapa: Imagen de Despliegue

#Utiliza una imagen base vacía, lo que minimiza el tamaño de la imagen final.
FROM scratch

#Define el directorio de trabajo.
WORKDIR /hipstershop

# Definition of this variable is used by 'skaffold debug' to identify a golang binary.
# Default behavior - a failure prints a stack trace for the current goroutine.
# See https://golang.org/pkg/runtime/

#Configura la variable de entorno GOTRACEBACK para controlar la información de depuración que se muestra en caso de fallo.
ENV GOTRACEBACK=single

COPY ./data /hipstershop/data/

# Copia los certificados SSL desde la etapa de construcción (builder) al contenedor final.
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

# Copia la herramienta grpc_health_probe desde la etapa de construcción.
COPY --from=builder /bin/grpc_health_probe /bin/

# Copia el binario del servicio construido en la primera etapa al contenedor final.
COPY --from=builder /go/src/hipstershop/productcatalogservice /hipstershop/productcatalogservice

# Define el binario productcatalogservice como el punto de entrada del contenedor. 
#   Este comando se ejecutará cuando el contenedor se inicie.
ENTRYPOINT ["/hipstershop/productcatalogservice"]