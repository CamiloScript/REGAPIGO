package handlers

import (
    "net/http"
    "github.com/CamiloScript/REGAPIGO/application/documento"
    "github.com/CamiloScript/REGAPIGO/shared/logger"
    "github.com/CamiloScript/REGAPIGO/shared/config"
    "github.com/gin-gonic/gin"
    "github.com/CamiloScript/REGAPIGO/infraestructure/persistence/db/mongo"
    "github.com/CamiloScript/REGAPIGO/domain/auth" 
    "github.com/CamiloScript/REGAPIGO/shared/utils"
)

// ManejadorDocumentos controla las operaciones con documentos.
type ManejadorDocumentos struct {
    servicio     *documento.ImplementacionServicioDocumentos // Servicio de documentos
    log          *logger.Registrador                         // Logger para registro de eventos
    cfg          *config.Config                              // Configuración de la aplicación
    internalAuth *InternalAuth                               // Autenticación interna
}

// NuevoManejadorDocumentos inicializa el manejador con dependencias.
func NuevoManejadorDocumentos(
    servicio *documento.ImplementacionServicioDocumentos,
    log *logger.Registrador,
    cfg *config.Config,
    authServicio auth.AuthService, // Servicio de autenticación
) *ManejadorDocumentos {
    return &ManejadorDocumentos{
        servicio:     servicio,
        log:          log,
        cfg:          cfg,
        internalAuth: NewInternalAuth(authServicio, log, cfg), // Inicializar autenticación interna
    }
}


// ManejadorSubirDocumento maneja la subida de documentos en formato base64.
func (h *ManejadorDocumentos) ManejadorSubirDocumento(c *gin.Context) {
// 1. Autenticación interna
ticket, err := h.internalAuth.AutenticarInternamente()
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno de autenticación"})
    return
}

// 2. Extraer archivo en base64 y metadatos del cuerpo de la solicitud
var solicitud struct {
    Base64    string                 `json:"base64"`    // Archivo en formato base64
    Metadatos map[string]interface{} `json:"metadatos"` // Metadatos en formato JSON
}
if err := c.ShouldBindJSON(&solicitud); err != nil {
    h.log.Error("Solicitud inválida", map[string]interface{}{"error": err.Error()})
    c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de solicitud incorrecto"})
    return
}

// 3. Decodificar el archivo base64 a bytes
fileBytes, err := utils.DecodeBase64(solicitud.Base64)
if err != nil {
    h.log.Error("Error al decodificar base64", map[string]interface{}{"error": err.Error()})
    c.JSON(http.StatusBadRequest, gin.H{"error": "Archivo base64 inválido"})
    return
}

// 4. Delegar al servicio de documentos
respuesta, err := h.servicio.SubirDocumento(c, fileBytes, solicitud.Metadatos, ticket)
if err != nil {
    h.log.Error("Error interno", map[string]interface{}{"error": err.Error()})
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar el documento"})
    return
}

// 5. Responder con éxito
c.JSON(http.StatusOK, respuesta)
h.log.Info("Documento subido", map[string]interface{}{"id": respuesta["entry"].(map[string]interface{})["id"]})

// 6. Persistir en MongoDB
if err := mongo.GuardarEnMongoDB(respuesta, h.log); err != nil {
    h.log.Error("Error en persistencia MongoDB", map[string]interface{}{"error": err.Error()})
}
}

// ManejadorListarDocumentos procesa el listado de documentos.
func (h *ManejadorDocumentos) ManejadorListarDocumentos(c *gin.Context) {
    // 1. Autenticación interna
    ticket, err := h.internalAuth.AutenticarInternamente()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno de autenticación"})
        return
    }

    // 2. Extraer filtros del cuerpo
    var filtros map[string]interface{}
    if err := c.ShouldBindJSON(&filtros); err != nil {
        h.log.Error("Error al analizar filtros", map[string]interface{}{"error": err.Error()})
        c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de filtros inválido"})
        return
    }

    // 3. Delegar al servicio de documentos
    documentos, err := h.servicio.ListarDocumentos(c, filtros, ticket)
    if err != nil {
        h.log.Error("Error al listar documentos", map[string]interface{}{"error": err.Error()})
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al listar documentos"})
        return
    }

    // 4. Responder con éxito
    c.JSON(http.StatusOK, gin.H{
        "data": documentos,
        "meta": map[string]interface{}{
            "total": len(documentos),
        },
    })
    h.log.Info("Documentos listados", map[string]interface{}{"total": len(documentos)})
}

// ManejadorDescargarDocumento maneja la descarga de documentos y los devuelve en formato base64.
func (h *ManejadorDocumentos) ManejadorDescargarDocumento(c *gin.Context) {
    // 1. Autenticación interna
    ticket, err := h.internalAuth.AutenticarInternamente()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno de autenticación"})
        return
    }

    // 2. Extraer ID del archivo desde parámetros de consulta
    idFile := c.Query("idFile")
    if idFile == "" {
        h.log.Error("ID de archivo faltante", nil)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere el parámetro idFile"})
        return
    }

    // 3. Delegar al servicio de documentos
    contenido, nombreArchivo, err := h.servicio.DescargarDocumento(c, idFile, ticket)
    if err != nil {
        if err == documento.ErrDocumentoNoEncontrado {
            h.log.Warn("Documento no encontrado", map[string]interface{}{"idFile": idFile})
            c.JSON(http.StatusNotFound, gin.H{"error": "documento no encontrado"})
        } else {
            h.log.Error("Error interno", map[string]interface{}{"error": err.Error()})
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al descargar el documento"})
        }
        return
    }

    // 4. Codificar el archivo a base64
    base64File := utils.EncodeToBase64(contenido)

    // 5. Configurar respuesta
    c.JSON(http.StatusOK, gin.H{
        "fileName": nombreArchivo,
        "base64":   base64File,
    })
    h.log.Info("Documento descargado", map[string]interface{}{"idFile": idFile, "nombreArchivo": nombreArchivo})
}