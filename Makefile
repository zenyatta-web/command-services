# Makefile
# El Makefile proporciona una manera conveniente de automatizar tareas comunes 
# relacionadas con el desarrollo, construcción, pruebas y despliegue del microservicio.

# Define una variable GOPATH que se utiliza para obtener la ruta del directorio Go workspace.
GOPATH:=$(shell go env GOPATH)

# Instala las dependencias y herramientas necesarias para el proyecto.
.PHONY: init
init:
	@go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

# Compila los archivos .proto en el proyecto utilizando el compilador Protobuf.
# Especifica la ruta donde se encuentran los archivos .proto.
# Especifica la salida para el generador de código Go y gRPC.
.PHONY: proto
proto:
	@protoc --proto_path=. --go_out=. --go-grpc_out=. proto/users.proto

# Actualiza las dependencias del proyecto.
.PHONY: update
update:
	@go get -u

# Compila el código Go en un ejecutable llamado productcatalogservice.
.PHONY: build
build:
	@go build -o productcatalogservice *.go

# Ejecuta las pruebas unitarias en el proyecto, mostrando la salida detallada y la cobertura del código.
.PHONY: test
test:
	@go test -v ./... -cover

# Construye una imagen Docker para el microservicio llamada productcatalogservice.
.PHONY: docker
docker:
	@docker build -t productcatalogservice:latest .
