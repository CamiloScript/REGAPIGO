package handlers

import (
    "encoding/json"
    "net/http"
    "fmt"
    "net/url"
    "github.com/CamiloScript/REGAPIGO/application/documento"
    "github.com/CamiloScript/REGAPIGO/domain/documentos"
    "github.com/CamiloScript/REGAPIGO/shared/logger"
    "github.com/CamiloScript/REGAPIGO/shared/config"
    "github.com/gin-gonic/gin"
    "github.com/CamiloScript/REGAPIGO/infraestructure/persistence/db/mongo"
    "github.com/CamiloScript/REGAPIGO/domain/auth" // Importar el paquete auth
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

// ManejadorSubirDocumento procesa la subida de documentos.
func (h *ManejadorDocumentos) ManejadorSubirDocumento(c *gin.Context) {
    // 1. Autenticación interna
    ticket, err := h.internalAuth.AutenticarInternamente()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno de autenticación"})
        return
    }

    // 2. Extraer archivo del formulario
    archivo, err := c.FormFile("documento")
    if err != nil {
        h.log.Error("Archivo no recibido", map[string]interface{}{"error": err.Error()})
        c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere un archivo"})
        return
    }

    // 3. Extraer metadatos del formulario
    metadatosStr := c.PostForm("propiedades")
    if metadatosStr == "" {
        h.log.Error("Metadatos faltantes", nil)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Metadatos requeridos"})
        return
    }

    // 4. Validar estructura de metadatos
    var metadatos documentos.DocumentMetadata
    if err := json.Unmarshal([]byte(metadatosStr), &metadatos); err != nil {
        h.log.Error("Metadatos inválidos", map[string]interface{}{"error": err.Error()})
        c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de metadatos incorrecto"})
        return
    }

    if err := metadatos.Validar(); err != nil {
        h.log.Error("Validación fallida", map[string]interface{}{"error": err.Error()})
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 4.1 Validar sub-categoría específica
    subCategoriasPermitidas := map[string]struct{}{
        "Riesgo":               {},
        "Comercial":            {},
        "Operaciones":         {},
        "Legal":               {},
        "Documentos de Cliente": {},
    }
    
    // Acceder al campo SubCategorias directamente desde el struct
    subCat := metadatos.SubCategorias
    if subCat == "" {
        h.log.Error("Campo 'tanner:sub-categorias' faltante", nil)
        c.JSON(http.StatusBadRequest, gin.H{"error": "El campo 'tanner:sub-categorias' es obligatorio"})
        return
    }

    if _, permitida := subCategoriasPermitidas[subCat]; !permitida {
        h.log.Error("Sub-categoría no permitida", map[string]interface{}{"valor": subCat})
        c.JSON(http.StatusBadRequest, gin.H{
            "error": fmt.Sprintf("Valor '%s' no permitido. Valores válidos: Riesgo, Comercial, Operaciones, Legal, Documentos de Cliente", subCat),
        })
        return
    }

    // 5. Delegar al servicio de documentos
    respuesta, err := h.servicio.SubirDocumento(c, archivo, metadatosStr, ticket)
    if err != nil {
        h.log.Error("Error interno", map[string]interface{}{"error": err.Error()})
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar el documento"})
        return
    }

    // 6. Responder con éxito
    c.JSON(http.StatusOK, respuesta)
    h.log.Info("Documento subido", map[string]interface{}{"id": respuesta["entry"].(map[string]interface{})["id"]})

    // 7. Persistir en MongoDB
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

// ManejadorDescargarDocumento procesa la descarga de documentos.
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

    // 4. Configurar headers mejorados para descarga de PDF
    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", url.QueryEscape(nombreArchivo)))
    c.Header("Content-Type", "application/pdf")
    c.Header("Content-Transfer-Encoding", "binary")
    c.Header("Expires", "0")
    c.Header("Cache-Control", "must-revalidate")
    c.Header("Pragma", "public")
    c.Header("Content-Length", fmt.Sprintf("%d", len(contenido)))

    // 5. Enviar el contenido del archivo al cliente
    c.Data(http.StatusOK, "application/pdf", contenido)
    h.log.Info("Documento descargado", map[string]interface{}{"idFile": idFile, "nombreArchivo": nombreArchivo})
}