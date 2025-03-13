package routes

import (
    "github.com/CamiloScript/REGAPIGO/domain/auth"
    "github.com/CamiloScript/REGAPIGO/infraestructure/api/handlers"
    "github.com/CamiloScript/REGAPIGO/infraestructure/persistence/servicio"
    "github.com/CamiloScript/REGAPIGO/shared/config"
    "github.com/CamiloScript/REGAPIGO/shared/logger"
    "github.com/CamiloScript/REGAPIGO/application/documento"
    "github.com/gin-gonic/gin"
)

// RegistrarRutas configura todas las rutas de la API.
// Parámetros:
//   - router: Instancia del enrutador de Gin.
//   - log: Logger para registrar eventos y errores.
//   - cfg: Configuración de la aplicación.
func RegistrarRutas(router *gin.Engine, log *logger.Registrador, cfg *config.Config) {

    // Autenticación
    clienteAuth := servicio.NewAuthClient(cfg, log)                 // Cliente de autenticación
    servicioAuth := auth.NewAuthService(clienteAuth, log)           // Servicio de autenticación
    manejadorAuth := handlers.NewAuthHandler(servicioAuth, log)     // Manejador de autenticación

    // Ruta de login
    router.POST("/auth/login", manejadorAuth.Login)

    // Grupo de rutas protegidas (documentos)
    grupoDocumentos := router.Group("/documentos")
    {
        // Inicializar servicios
        servicioDocs := servicio.NuevoServicioDocumentos(cfg, log)                      // Servicio de Alfresco
        servicioNegocio := documento.NuevoServicioDocumentos(servicioDocs, log, cfg.AlfrescoAPIKey) // Servicio de negocio

        // Inicializar manejadores con autenticación interna
        manejadorDocs := handlers.NuevoManejadorDocumentos(servicioNegocio, log, cfg, servicioAuth) // Manejador de documentos
        manejadorBusqueda := handlers.NuevoManejadorBusquedaDescarga(servicioNegocio, log, cfg, servicioAuth) // Manejador de búsqueda y descarga

        // Ruta para subir documentos: recibe una solicitud POST en "/documentos/subir".
        grupoDocumentos.POST("/subir", manejadorDocs.ManejadorSubirDocumento)

        // Ruta para listar documentos: recibe una solicitud POST en "/documentos/listar".
        grupoDocumentos.POST("/listar", manejadorDocs.ManejadorListarDocumentos)

        // Ruta para descargar documentos: recibe una solicitud POST en "/documentos/descargar".
        grupoDocumentos.POST("/descargar", manejadorDocs.ManejadorDescargarDocumento)
        
        // Ruta para subir lotes de documentos: recibe una solicitud POST en "/documentos/subir-lote".
        grupoDocumentos.POST("/subir-lote", manejadorDocs.ManejadorLoteDocumentos)

        // Ruta para buscar y descargar documentos: recibe una solicitud POST en "/documentos/buscar-descargar".
        grupoDocumentos.POST("/buscar-descargar", manejadorBusqueda.BuscarYDescargarDocumento)
    }

    // Ruta de health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"estado": "activo"})
    })
}