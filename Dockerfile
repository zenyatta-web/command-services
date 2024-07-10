# Este Dockerfile está diseñado para crear una imagen Docker para un microservicio escrito en Go.

# 1. Primera Etapa: Builder
FROM golang:1.20.2-alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux

WORKDIR /go/src/zen-command-services

RUN apk --update --no-cache add ca-certificates make protoc

COPY Makefile go.mod go.sum ./

RUN go env -w GOPROXY=https://goproxy.io,direct/

RUN make init && go mod download

COPY . .

RUN go build -o /go/src/zen-command-services/zen-command-services .

# 2. Segunda Etapa: Imagen de Despliegue
FROM scratch

WORKDIR /zen-command-services

# Copiar el archivo .env al contenedor
COPY .env /zen-command-services/.env

COPY ./data /zen-command-services/data/
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /go/src/zen-command-services/zen-command-services /zen-command-services/zen-command-services

# Definir variables de entorno (opcional si las cargas desde .env)
ENV PORT=50052
ENV MONGO_URI="mongodb+srv://mendieta19ns:A1em3yKnw7RmqvL0@zen.xqsdoat.mongodb.net/"
ENV MONGO_DATABASE="zen"

ENTRYPOINT ["/zen-command-services/zen-command-services"]
