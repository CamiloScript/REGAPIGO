package main

import (
    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "github.com/CamiloScript/REGAPIGO/src/APIGuest/APIController"
    "github.com/CamiloScript/REGAPIGO/bd/MongoDB"
    "os"
    "time"
)

func main() {
    // Se configura zerolog en un formato de tiempo específico (RFC3339) y se establece la salida a la consola.
    zerolog.TimeFieldFormat = time.RFC3339
    // Se establece el logger de zerolog para que escriba en la salida estándar.
    log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

    // Se inicializa el servidor de Gin, ingresandolo a la variable router.
    router := gin.New()
    
    // Middlewares (conservando Recovery y agregando nuestro logger) 
    router.Use(ginLogger())  // Logger personalizado con zerolog
    router.Use(gin.Recovery())  // Mantenemos el Recovery de Gin para manejo de errores en tiempo de ejecución.

    // Conexión a la base de datos 
    _, err := MongoDB.ConexionDB()
    // Si hay un error al conectar a la base de datos, se registra en el log y se finaliza la ejecución.
    if err != nil {
        log.Fatal().
            Err(err).
            Msg("Error al conectar a la base de datos")
    }
    defer MongoDB.CerrarConexion()

    // Definimos las rutas de la API para controlar las solicitudes HTTP.
    router.GET("/guests", APIController.GetGuests)
    router.GET("/guest/:id", APIController.GetGuestByID)
    router.POST("/guests", APIController.CreateGuest)
    router.PUT("/guest/:id", APIController.UpdateGuest)
    router.DELETE("/guest/:id", APIController.DeleteGuest)

    // Iniciamos el servidor en el puerto 8080 y se registra en el log.
    log.Info().Msg("Iniciando servidor en puerto :8080")
    if err := router.Run(":8080"); err != nil {
        log.Fatal().
            Err(err).
            Msg("Error al iniciar el servidor")
    }
}

// Se genera la función ginLogger(), la que retorna un middleware para registrar las solicitudes HTTP con zerolog.
func ginLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Se registra el tiempo de inicio de la solicitud.
        start := time.Now()
        // Se ejecuta el siguiente middleware.
        c.Next()
        // Se registra en el log la solicitud HTTP, con información como el método, la ruta, el estado, la IP y la duración.
        // Se busca que el log de la solicitud http, sea personalizado y con la información necesaria para su seguimiento.
        log.Info().
            Str("método", c.Request.Method).
            Str("ruta", c.Request.URL.Path).
            Int("status", c.Writer.Status()).
            Str("ip", c.ClientIP()).
            Dur("duración", time.Since(start)).
            Msg("solicitud HTTP")
    }
}