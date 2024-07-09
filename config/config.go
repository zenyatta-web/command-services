/*
Este archivo de configuración en Go utiliza el paquete go-micro para cargar y gestionar
la configuración de tu microservicio, permitiendo flexibilidad y facilidad de uso al
soportar variables de entorno.
*/
package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Port    int
	Tracing TraicingConfig //rastreo
	Mongo   MongoConfig
}

type TraicingConfig struct {
	Enable bool
	Jaeger JaegerConfig
}

type JaegerConfig struct {
	URL string
}

type MongoConfig struct {
	URI      string
	Database string
}

var cfg *Config = &Config{
	Port: 50052,
}

func Address() string {
	return fmt.Sprintf(":%d", cfg.Port)
}

func Tracing() TraicingConfig {
	return cfg.Tracing
}

func Mongo() MongoConfig {
	return cfg.Mongo
}

func Load() error {
	// Cargar variables de entorno desde un archivo .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error cargando las variables de entorno: ", err)
	}

	// Configurar viper para leer variables de entorno
	viper.AutomaticEnv()

	cfg.Port = viper.GetInt("PORT")
	cfg.Mongo.URI = viper.GetString("MONGO_URI")
	cfg.Mongo.Database = viper.GetString("MONGO_DATABASE")

	return nil
}
