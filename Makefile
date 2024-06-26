# El Makefile proporciona una manera conveniente de automatizar tareas comunes 
# 	relacionadas con el desarrollo, construcción, pruebas 
# 	 y despliegue del microservicio.


# Define una variable GOPATH que se utiliza para obtener la ruta del directorio Go workspace.
GOPATH:=$(shell go env GOPATH)

# Instala las dependencias y herramientas necesarias para el proyecto.
# Descarga y actualiza la librería protobuf de Google.
# Instala los generadores de código protobuf para Go y Micro.
.PHONY: init
init:
	@go get -u google.golang.org/protobuf/proto
	@go install github.com/golang/protobuf/protoc-gen-go@latest
	@go install github.com/go-micro/generator/cmd/protoc-gen-micro@latest


# Compila los archivos .proto en el proyecto utilizando el compilador Protobuf.
# Especifica la ruta donde se encuentran los archivos .proto.
# Especifica la salida para el generador de código Micro.
# Especifica la salida para el generador de código Go.
.PHONY: proto
proto:
	@protoc --proto_path=. --micro_out=. --go_out=:. proto/productcatalogservice.proto
	@protoc --proto_path=. --micro_out=. --go_out=:. proto/health.proto
	

# Actualiza las dependencias del proyecto.
.PHONY: update
update:
	@go get -u

# Realiza una limpieza de los módulos Go y sus dependencias.
# .PHONY: tidy
# tidy:
# 	@go mod tidy

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