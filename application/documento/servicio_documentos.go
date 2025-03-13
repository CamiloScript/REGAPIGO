package documento

import (
    "github.com/gin-gonic/gin"
    "github.com/CamiloScript/REGAPIGO/domain/documentos"
    "github.com/CamiloScript/REGAPIGO/shared/logger"
    "mime/multipart"
    "fmt"
    "errors"
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

// SubirDocumento maneja la subida de un documento, delegando la operación al almacenamiento.
func (s *ImplementacionServicioDocumentos) SubirDocumento(
    c *gin.Context,                  // Contexto de Gin para manejar la solicitud HTTP
    archivo *multipart.FileHeader,   // Archivo a subir
    metadatos string,                // Metadatos asociados al archivo
    ticket string,                   // Ticket obtenido del cliente para autenticación/autorización
) (map[string]interface{}, error) {  // Retorna un mapa con la respuesta o un error

    // Delegar la operación de subida al almacenamiento
    respuesta, err := s.almacenamiento.SubirDocumento(c, archivo, metadatos, ticket, s.apiKey)
    if err != nil {
        // Registrar el error en el logger
        s.log.Error("Error en el servicio", map[string]interface{}{"error": err.Error()})
        // Retornar un error formateado para el cliente
        return nil, fmt.Errorf("error interno: %v", err)
    }

    // Retornar la respuesta del almacenamiento
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