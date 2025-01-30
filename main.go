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
    // Configurar zerolog con formato legible
    zerolog.TimeFieldFormat = time.RFC3339
    log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

    // Inicializamos el servidor de Gin (manteniendo tu variable 'router')
    router := gin.New()
    
    // Middlewares (conservando Recovery y agregando nuestro logger)
    router.Use(ginLogger())  // Logger personalizado
    router.Use(gin.Recovery())  // Mantenemos el Recovery de Gin

    // Conexión a la base de datos (mismo código con logging mejorado)
    _, err := MongoDB.ConexionDB()
    if err != nil {
        log.Fatal().
            Err(err).
            Msg("Error al conectar a la base de datos")
    }
    defer MongoDB.CerrarConexion()

    // Definimos las rutas (sin cambios)
    router.GET("/guests", APIController.GetGuests)
    router.GET("/guest/:id", APIController.GetGuestByID)
    router.POST("/guests", APIController.CreateGuest)
    router.PUT("/guest/:id", APIController.UpdateGuest)
    router.DELETE("/guest/:id", APIController.DeleteGuest)

    // Iniciamos el servidor (logging mejorado)
    log.Info().Msg("Iniciando servidor en puerto :8080")
    if err := router.Run(":8080"); err != nil {
        log.Fatal().
            Err(err).
            Msg("Error al iniciar el servidor")
    }
}

// ginLogger es un middleware personalizado con zerolog
func ginLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        
        log.Info().
            Str("método", c.Request.Method).
            Str("ruta", c.Request.URL.Path).
            Int("status", c.Writer.Status()).
            Str("ip", c.ClientIP()).
            Dur("duración", time.Since(start)).
            Msg("solicitud HTTP")
    }
}