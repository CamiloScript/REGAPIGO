package MongoDB

import (
	"context"
	"fmt"
	"log"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Cliente *mongo.Client

// ConexionDB establece la conexión con la base de datos de MongoDB
func ConexionDB() (*mongo.Client, error) {
	// Establecemos un contexto con un timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Cadena de conexión a MongoDB
	uri := "mongodb://localhost:27017"

	// Intentamos conectar al servidor de MongoDB
	Cliente, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Println("Error al conectar con MongoDB:", err)
		return nil, fmt.Errorf("error al conectar con MongoDB: %v", err)
	}

	// Verificamos la conexión
	if err := Cliente.Ping(ctx, nil); err != nil {
		log.Println("Error al hacer ping a MongoDB:", err)
		return nil, fmt.Errorf("error al hacer ping a MongoDB: %v", err)
	}

	log.Println("Conexión exitosa a MongoDB")
	return Cliente, nil
}

// CerrarConexion cierra la conexión a la base de datos
func CerrarConexion() {
	if err := Cliente.Disconnect(context.Background()); err != nil {
		log.Fatal("Error al desconectar de MongoDB:", err)
	}
	log.Println("Conexión cerrada con éxito")
}
