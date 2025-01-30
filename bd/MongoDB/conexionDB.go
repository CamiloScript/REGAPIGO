package MongoDB

import (
    "context"
    "fmt"
    "os"
    "time"
	
    "github.com/joho/godotenv"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var Cliente *mongo.Client

// ConexionDB establece la conexión con MongoDB usando zerolog
func ConexionDB() (*mongo.Client, error) {
    // Configurar zerolog (formato legible para humanos)
    log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

    // Cargar variables de entorno
    err := godotenv.Load()
    if err != nil {
        log.Fatal().
            Err(err).
            Msg("Error al cargar el archivo .env")
        return nil, fmt.Errorf("error al cargar el archivo .env: %v", err)
    }

    // Leer variable de entorno
    uri := os.Getenv("MONGO_URI")
    if uri == "" {
        log.Fatal().
            Str("variable", "MONGO_URI").
            Msg("Variable de entorno no definida")
        return nil, fmt.Errorf("la variable de entorno MONGO_URI no está definida")
    }

    // Contexto con timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Conectar a MongoDB
    Cliente, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        log.Error().
            Err(err).
            Str("uri", uri).
            Msg("Error de conexión a MongoDB")
        return nil, fmt.Errorf("error al conectar con MongoDB: %v", err)
    }

    // Verificar conexión
    if err := Cliente.Ping(ctx, nil); err != nil {
        log.Error().
            Err(err).
            Msg("Error verificando conexión")
        return nil, fmt.Errorf("error al hacer ping a MongoDB: %v", err)
    }

    log.Info().
        Msg("Conexión exitosa a MongoDB")
    
    return Cliente, nil
}

// CerrarConexion cierra la conexión con MongoDB usando zerolog
func CerrarConexion() {
    if Cliente != nil {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        if err := Cliente.Disconnect(ctx); err != nil {
            log.Error().
                Err(err).
                Msg("Error al desconectar de MongoDB")
        } else {
            log.Info().
                Msg("Conexión cerrada con éxito")
        }
    }
}
