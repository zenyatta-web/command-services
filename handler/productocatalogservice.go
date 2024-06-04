package handler

/*
	Este archivo define el servicio ProductCatalogService, que maneja la lógica de negocio
	relacionada con la lista de productos en el catálogo, proporcionando varias
	funcionalidades a través de métodos gRPC.
*/

import (
	"context"
	"io/ioutil" //Para leer archivos.
	"os"        //Para manejar señales del sistema operativo.
	"os/signal" //Para manejar señales del sistema operativo.
	"strings"
	"sync"    //Para sincronización de goroutines.
	"syscall" //Para definir señales del sistema operativo.

	"go-micro.dev/v4/logger"                        //Para registrar mensajes de log.
	"google.golang.org/grpc/codes"                  //Para manejar códigos de estado y errores en gRPC.
	"google.golang.org/grpc/status"                 //Para manejar códigos de estado y errores en gRPC.
	"google.golang.org/protobuf/encoding/protojson" //Para trabajar con JSON y Protobuf.

	/*
		El paquete generado a partir de las definiciones de Protobuf, contiene los mensajes y servicios
		definidos en los archivos .proto.
	*/
	pb "github.com/zenyatta-web/command-services/proto"
)

var reloadCatalog bool // Una bandera para determinar si el catálogo debe ser recargado.

type ProductCatalogService struct {
	sync.Mutex //Para asegurar que las operaciones de lectura/escritura en products sean seguras en concurrencia.
	//Un slice de punteros a pb.Product, que almacena la lista de productos cargados del archivo JSON.
	products []*pb.Product
}

/*
Llena la respuesta con todos los productos del catálogo. Llama a s.parseCatalog para obtener
la lista actualizada de productos.
*/
func (s *ProductCatalogService) ListProducts(ctx context.Context, in *pb.Empty, out *pb.ListProductsResponse) error {
	out.Products = s.parseCatalog()
	return nil
}

/*
Busca un producto por ID en el catálogo y llena la respuesta con los detalles del producto
encontrado. Si no se encuentra, devuelve un error NotFound.
*/
func (s *ProductCatalogService) GetProduct(ctx context.Context, in *pb.GetProductRequest, out *pb.Product) error {
	var found *pb.Product
	products := s.parseCatalog()
	for _, p := range products {
		if in.Id == p.Id {
			found = p
		}
	}
	if found == nil {
		return status.Errorf(codes.NotFound, "no product with ID %s", in.Id)
	}
	out.Id = found.Id
	out.Name = found.Name
	out.Categories = found.Categories
	out.Description = found.Description
	out.Picture = found.Picture
	out.PriceUsd = found.PriceUsd
	return nil
}

/*
Busca productos cuyo nombre o descripción contenga la cadena de consulta (query).
Llena la respuesta con los productos que coinciden.
*/
func (s *ProductCatalogService) SearchProducts(ctx context.Context, in *pb.SearchProductsRequest, out *pb.SearchProductsResponse) error {
	var ps []*pb.Product
	products := s.parseCatalog()
	for _, p := range products {
		if strings.Contains(strings.ToLower(p.Name), strings.ToLower(in.Query)) ||
			strings.Contains(strings.ToLower(p.Description), strings.ToLower(in.Query)) {
			ps = append(ps, p)
		}
	}
	out.Results = ps
	return nil
}

/*
Lee el archivo JSON del catálogo de productos y lo deserializa en una estructura ListProductsResponse.
Usa un bloqueo para asegurar la operación en concurrencia.
*/
func (s *ProductCatalogService) readCatalogFile() (*pb.ListProductsResponse, error) {
	s.Lock()
	defer s.Unlock()
	catalogJSON, err := ioutil.ReadFile("data/products.json")
	if err != nil {
		logger.Errorf("failed to open product catalog json file: %v", err)
		return nil, err
	}
	catalog := &pb.ListProductsResponse{}
	if err := protojson.Unmarshal(catalogJSON, catalog); err != nil {
		logger.Warnf("failed to parse the catalog JSON: %v", err)
		return nil, err
	}
	logger.Info("successfully parsed product catalog json")
	return catalog, nil
}

/*
Verifica si debe recargar el catálogo (reloadCatalog es verdadero o products está vacío).
Si es necesario, lee el archivo del catálogo. Devuelve la lista de productos.
*/
func (s *ProductCatalogService) parseCatalog() []*pb.Product {
	if reloadCatalog || len(s.products) == 0 {
		catalog, err := s.readCatalogFile()
		if err != nil {
			return []*pb.Product{}
		}
		s.products = catalog.Products
	}
	return s.products
}

/*
Configura un manejador de señales para SIGUSR1 y SIGUSR2. Dependiendo de la señal recibida,
establece la bandera reloadCatalog para habilitar o deshabilitar la recarga del catálogo.
*/
func init() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for {
			sig := <-sigs
			logger.Infof("Received signal: %s", sig)
			if sig == syscall.SIGUSR1 {
				reloadCatalog = true
				logger.Infof("Enable catalog reloading")
			} else {
				reloadCatalog = false
				logger.Infof("Disable catalog reloading")
			}
		}
	}()
}
