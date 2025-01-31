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

// La función ConexionDB() establece la conexión con MongoDB usando zerolog para registrar los eventos.
func ConexionDB() (*mongo.Client, error) {
    // Se configura zerolog (formato predeterminado)
    log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

    // Se cargan las variables de entorno, presentes en el archivo .env, para la conexion a mongoDB
    err := godotenv.Load()
    // Se genera un error si no se puede cargar el archivo .env, y se registra con zerolog para su seguimiento.
    if err != nil {
        log.Fatal().
            Err(err).
            Msg("Error al cargar el archivo .env")
        return nil, fmt.Errorf("error al cargar el archivo .env: %v", err)
    }

    // Se lee la variable de entorno MONGO_URI, presente dentro del archivo .env
    uri := os.Getenv("MONGO_URI")
    // Se genera un error si la variable de entorno MONGO_URI no está definida, y se registra con zerolog para su seguimiento.
    if uri == "" {
        log.Fatal().
            Str("variable", "MONGO_URI").
            Msg("Variable de entorno no definida")
        return nil, fmt.Errorf("la variable de entorno MONGO_URI no está definida")
    }

    // Contexto con timeout de 10 segundos para la conexión a MongoDB, y se cancela al finalizar la función.
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Se genera la conexion a mongoDB, y se registra con zerolog para su seguimiento.
    Cliente, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        log.Error().
            Err(err).
            Str("uri", uri).
            Msg("Error de conexión a MongoDB")
        return nil, fmt.Errorf("error al conectar con MongoDB: %v", err)
    }

    // Se verifica la conexion a mongoDB a traves de un ping, un ping es una solicitud de prueba para verificar la conexión con el servidor.
    // Se registran los eventos con zerolog para su seguimiento.
    if err := Cliente.Ping(ctx, nil); err != nil {
        log.Error().
            Err(err).
            Msg("Error verificando conexión")
        return nil, fmt.Errorf("error al hacer ping a MongoDB: %v", err)
    }
    // Se registra con zerolog la conexión exitosa a MongoDB.
    log.Info().
        Msg("Conexión exitosa a MongoDB")
    
    return Cliente, nil
}

// Se cierra la conexion a mongo db cuando se llame a la funcion CerrarConexion(), la cual se encarga de cerrar la conexion con mongoDB.
func CerrarConexion() {
    // Se verifica si la conexión a MongoDB es diferente de nil, para cerrar la conexión.
    if Cliente != nil {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        // Se cierra la conexión a MongoDB, y se registra con zerolog para su seguimiento.
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
