# Define una variable GOPATH que se utiliza para obtener la ruta del directorio Go workspace.
GOPATH:=$(shell go env GOPATH)

# Instala las dependencias y herramientas necesarias para el proyecto.
# Descarga y actualiza la librería protobuf de Google.
# Instala los generadores de código protobuf para Go y gRPC.
.PHONY: init
init:
	@go get -u google.golang.org/protobuf/proto
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

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

# Realiza una limpieza de los módulos Go y sus dependencias.
.PHONY: tidy
tidy:
	@go mod tidy

# Compila el código Go en un ejecutable llamado zen-command-services.
.PHONY: build
build:
	@go build -o zen-command-services *.go

# Ejecuta las pruebas unitarias en el proyecto, mostrando la salida detallada y la cobertura del código.
.PHONY: test
test:
	@go test -v ./... -cover

# Construye una imagen Docker para el microservicio llamada zen-command-services.
.PHONY: docker
docker:
	@docker build -t zen-command-services:latest .
