package MongoDB

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/joho/godotenv"
)

// ConexionDB establece la conexión con MongoDB utilizando URI desde el archivo .env.
func ConexionDB() (*mongo.Client, error) {
	// Cargar variables del archivo .env
	err := godotenv.Load("/")
	if err != nil {
		log.Error().Err(err).Msg("Error al cargar el archivo .env")
		return nil, fmt.Errorf("error al cargar el archivo .env: %v", err)
	}

	// Obtener URI de la base de datos desde la variable de entorno
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Error().Msg("URI de MongoDB no encontrada en las variables de entorno")
		return nil, fmt.Errorf("URI de MongoDB no encontrada en las variables de entorno")
	}
	log.Debug().Msg("Conectando a MongoDB con URI: " + uri)  // Esto te permitirá ver la URI en los logs

	// Opciones de conexión
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Crear contexto con timeout de 10 segundos
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Crear cliente y conectarse a la base de datos
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Error().Err(err).Msg("Error al conectar con la base de datos MongoDB")
		return nil, fmt.Errorf("error al conectar con MongoDB: %v", err)
	}

	// Verificar conexión mediante un ping
	err = client.Database("admin").RunCommand(ctx, bson.D{{"ping", 1}}).Err()
	if err != nil {
		log.Error().Err(err).Msg("Error al hacer ping a MongoDB")
		return nil, fmt.Errorf("error al hacer ping a MongoDB: %v", err)
	}

	// Registro de conexión exitosa
	log.Info().Msg("Conexión exitosa a la base de datos MongoDB")

	return client, nil
}

// CerrarConexion cierra la conexión con MongoDB.
func CerrarConexion(client *mongo.Client) {
	// Crear contexto con timeout para la desconexión
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Intentar desconectar el cliente
	err := client.Disconnect(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error al desconectar de MongoDB")
	} else {
		log.Info().Msg("Desconexión exitosa de MongoDB")
	}
}

// ObtenerBaseDatos devuelve una base de datos específica
func ObtenerBaseDatos(client *mongo.Client, nombreBD string) *mongo.Database {
    return client.Database(nombreBD)
}