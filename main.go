package main

import (
    "github.com/CamiloScript/REGAPIGO/shared/middleware"
    "github.com/CamiloScript/REGAPIGO/infraestructure/routes"
    "github.com/CamiloScript/REGAPIGO/shared/config"
    "github.com/CamiloScript/REGAPIGO/shared/logger"
    "github.com/CamiloScript/REGAPIGO/infraestructure/persistence/db/mongo"
    "github.com/gin-gonic/gin"
    ginSwagger "github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"
    "github.com/gin-contrib/cors"
)



func main() {
    // 1. Cargar configuración desde variables de entorno
    cfg := config.CargarConfiguracion()

    // 2. Validar variables críticas antes de continuar
    if cfg.MongoURI == "" {
        panic("Falta configurar MONGODB_URI en el archivo .env")
    }

    // 3. Inicializar logger con formato unificado
    log := logger.NuevoRegistrador("API_DOC", cfg.SeparadorLog)
    log.Info("Inicializando servicios", map[string]interface{}{"puerto": cfg.Puerto})

    // 4. Conectar a MongoDB (verificación temprana)
    if _, err := mongo.ConexionDB(); err != nil {
        log.Fatal("Error crítico en MongoDB", map[string]interface{}{
            "modulo": "database",
            "error":  err.Error(),
        })
    }
    log.Info("Conexión MongoDB establecida", nil)

    

    // 5. Configurar router Gin con middlewares
    router := gin.Default()
    router.Use(
        middleware.MiddlewareRegistro(log), // Middleware de logging
    )

    // Configurar CORS
    router.Use(cors.New(cors.Config{
        AllowAllOrigins: true,
        AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:    []string{"Origin", "Content-Type", "Accept", "Authorization", "ADFTannerServices"},
    }))

    // 6. Registrar todas las rutas HTTP
    routes.RegistrarRutas(router, log, cfg)

    // 7. Servir archivos estáticos
    router.Static("/docs", "./docs")

    // 8. Registrar ruta de Swagger
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // 9. Iniciar servidor HTTP
    log.Info("Servidor listo", map[string]interface{}{
        "puerto":  cfg.Puerto,
        "version": "1.2.0",
    })
    if err := router.Run(":" + cfg.Puerto); err != nil {
        log.Fatal("Error al iniciar servidor", map[string]interface{}{"error": err.Error()})
    }
}