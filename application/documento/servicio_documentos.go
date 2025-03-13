package documento

import (
    "github.com/gin-gonic/gin"
    "github.com/CamiloScript/REGAPIGO/domain/documentos"
    "github.com/CamiloScript/REGAPIGO/shared/utils"
    "github.com/CamiloScript/REGAPIGO/shared/logger"
    "mime/multipart"
    "fmt"
    "errors"
    "io"
    "encoding/json"
)

// ErrDocumentoNoEncontrado es un error que se produce cuando un documento no se encuentra.
var ErrDocumentoNoEncontrado = errors.New("documento no encontrado")

// ImplementacionServicioDocumentos maneja la lógica de documentos.
type ImplementacionServicioDocumentos struct {
    almacenamiento documentos.AlmacenamientoDocumentos // Almacenamiento inyectado para manejar la persistencia de documentos
    log            *logger.Registrador               // Logger para registrar eventos y errores
    apiKey         string                            // API Key desde configuración para autenticación
}

// NuevoServicioDocumentos construye el servicio con dependencias.
func NuevoServicioDocumentos(
    almacenamiento documentos.AlmacenamientoDocumentos, // Almacenamiento de documentos
    log *logger.Registrador,                           // Logger para registro de eventos
    apiKey string,                                     // API Key para autenticación
) *ImplementacionServicioDocumentos {
    return &ImplementacionServicioDocumentos{
        almacenamiento: almacenamiento, // Inicializa el almacenamiento
        log:            log,            // Inicializa el logger
        apiKey:         apiKey,         // Inicializa la API Key
    }
}

// SubirDocumento maneja la subida de un documento, validando el tipo de archivo.
func (s *ImplementacionServicioDocumentos) SubirDocumento(
    c *gin.Context,
    archivo *multipart.FileHeader,
    metadatos map[string]interface{},
    ticket string,
) (map[string]interface{}, error) {
    // 1. Leer el contenido del archivo
    file, err := archivo.Open()
    if err != nil {
        return nil, fmt.Errorf("no se pudo abrir el archivo: %v", err)
    }
    defer file.Close()

    // Leer los primeros bytes del archivo para detectar el tipo MIME
    buffer := make([]byte, 512)
    if _, err := file.Read(buffer); err != nil && err != io.EOF {
        return nil, fmt.Errorf("no se pudo leer el archivo: %v", err)
    }

    // 2. Validar el tipo de archivo
    mimeType := utils.DetectMimeTypeFromContent(buffer)
    if mimeType != "image/jpeg" && mimeType != "image/png" && mimeType != "application/pdf" {
        return nil, fmt.Errorf("tipo de archivo no soportado: %s", mimeType)
    }

    // 3. Convertir metadatos a JSON
    metadatosJSON, err := json.Marshal(metadatos)
    if err != nil {
        return nil, fmt.Errorf("error al convertir metadatos a JSON: %v", err)
    }

    // 4. Delegar la operación de subida al almacenamiento
    respuesta, err := s.almacenamiento.SubirDocumento(c, archivo, string(metadatosJSON), ticket, s.apiKey)
    if err != nil {
        s.log.Error("Error en el servicio", map[string]interface{}{"error": err.Error()})
        return nil, fmt.Errorf("error interno: %v", err)
    }

    return respuesta, nil
}

// ListarDocumentos maneja la lista de documentos, delegando la operación al almacenamiento.
func (s *ImplementacionServicioDocumentos) ListarDocumentos(
    c *gin.Context,                  // Contexto de Gin para manejar la solicitud HTTP
    filtros map[string]interface{},  // Filtros para la búsqueda de documentos
    ticket string,                   // Ticket obtenido del cliente para autenticación/autorización
) (map[string]interface{}, error) {  // Retorna un mapa con la respuesta o un error

    // Delegar la operación de listado al almacenamiento
    return s.almacenamiento.ListarDocumentos(c, filtros, ticket, s.apiKey)
}

// DescargarDocumento maneja la descarga de un documento, delegando la operación al almacenamiento.
func (s *ImplementacionServicioDocumentos) DescargarDocumento(
    c *gin.Context,                  // Contexto de Gin para manejar la solicitud HTTP
    idFile string,                   // ID del archivo a descargar
    ticket string,                   // Ticket obtenido del cliente para autenticación/autorización
) ([]byte, string, error) {          // Retorna el contenido del archivo, el tipo MIME y un error si lo hubiera

    // Delegar la operación de descarga al almacenamiento
    return s.almacenamiento.DescargarDocumento(c, idFile, ticket, s.apiKey)
}