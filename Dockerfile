# Dockerfile
# Este Dockerfile está diseñado para crear una imagen Docker para un microservicio escrito en Go.

# 1. Primera Etapa: Builder
# Utiliza la imagen base de Golang 1.20.2 sobre Alpine Linux, lo cual es una versión ligera de Linux. 
# La etiqueta AS builder define un alias para esta etapa de construcción.
FROM golang:1.20.2-alpine AS builder

# Establece variables de entorno para deshabilitar CGO (código C en Go) y especificar 
# que el sistema operativo objetivo es Linux.
ENV CGO_ENABLED=0 GOOS=linux

# Define el directorio de trabajo dentro del contenedor.
WORKDIR /go/src/zen-command-services

# Utiliza el gestor de paquetes apk para instalar certificados SSL, 
# la herramienta make y protoc (compilador de Protocol Buffers).
RUN apk --update --no-cache add ca-certificates make protoc

# Copia el Makefile y los archivos de dependencias de Go al contenedor.
COPY Makefile go.mod go.sum ./

# Configura el proxy de Go.
RUN go env -w GOPROXY=https://goproxy.io,direct/

# Ejecuta comandos para inicializar el entorno de desarrollo y descargar las dependencias de Go.
RUN make init && go mod download

# Copia todos los archivos del proyecto al contenedor.
COPY . .

# Compila el binario de Go.
RUN go build -o /go/src/zen-command-services/productcatalogservice .

# 2. Segunda Etapa: Imagen de Despliegue
FROM scratch

# Define el directorio de trabajo.
WORKDIR /zen-command-services

# Configura la variable de entorno GOTRACEBACK para controlar la información de depuración que se muestra en caso de fallo.
ENV GOTRACEBACK=single

COPY ./data /zen-command-services/data/

# Copia los certificados SSL desde la etapa de construcción (builder) al contenedor final.
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

# Copia el binario del servicio construido en la primera etapa al contenedor final.
COPY --from=builder /go/src/zen-command-services/productcatalogservice /zen-command-services/productcatalogservice

# Define el binario productcatalogservice como el punto de entrada del contenedor. 
ENTRYPOINT ["/zen-command-services/productcatalogservice"]
