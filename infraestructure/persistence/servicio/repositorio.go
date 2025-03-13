package servicio

import (
    "github.com/gin-gonic/gin"
    "github.com/CamiloScript/REGAPIGO/domain/documentos"
    "github.com/CamiloScript/REGAPIGO/shared/logger"
    "github.com/CamiloScript/REGAPIGO/shared/config"
)

// ServicioAlfresco implementa AlmacenamientoDocumentos para Alfresco.
type ServicioAlfresco struct {
    urlBase string                 // URL base del servidor Alfresco
    apiKey  string                 // API Key para autenticación
    cliente *ClienteAlfresco       // Cliente HTTP inyectado
    log     *logger.Registrador    // Logger para registrar eventos y errores
}

// NuevoServicioDocumentos inicializa el servicio.
// Ahora incluye configuración completa para integración con Alfresco.
// Parámetros:
//   - cfg: Configuración de la aplicación.
//   - log: Logger para registrar eventos y errores.
// Retorna una instancia configurada de ServicioAlfresco.
func NuevoServicioDocumentos(
    cfg *config.Config, 
    log *logger.Registrador,
) documentos.AlmacenamientoDocumentos {
    
    // Crear cliente HTTP configurado para Alfresco
    clienteAlfresco := NuevoClienteAlfresco(
        cfg.AlfrescoBaseURL,
        cfg.AlfrescoAPIKey, // Inyectar API Key desde configuración
        log,
    )

    return &ServicioAlfresco{
        urlBase: cfg.AlfrescoBaseURL,
        apiKey:  cfg.AlfrescoAPIKey,  // Nueva propiedad requerida
        cliente: clienteAlfresco,     // Cliente HTTP inyectado
        log:     log,
    }
}


// SubirDocumento implementa la subida stateless.
// Parámetros:
//   - c: Contexto de Gin.
//   - fileBytes: Bytes del archivo a subir.
//   - metadatos: Metadatos del documento en formato JSON.
//   - ticket: Ticket de autenticación.
//   - apiKey: API Key para autenticación.
// Retorna un mapa con la respuesta de Alfresco o un error en caso de fallo.
func (s *ServicioAlfresco) SubirDocumento(
    c *gin.Context,
    fileBytes []byte,
    metadatos string,
    ticket string,
    apiKey string,
) (map[string]interface{}, error) {

    // Crear cliente con API Key y ticket actual
    cliente := NuevoClienteAlfresco(s.urlBase, apiKey, s.log)
    return cliente.SubirDocumento(c.Request.Context(), fileBytes, metadatos, ticket)
}

// ListarDocumentos implementa el listado stateless.
// Parámetros:
//   - c: Contexto de Gin.
//   - filtros: Filtros de búsqueda en formato JSON.
//   - ticket: Ticket de autenticación.
//   - apiKey: API Key para autenticación.
// Retorna un mapa con los documentos listados o un error en caso de fallo.
func (s *ServicioAlfresco) ListarDocumentos(
    c *gin.Context,
    filtros map[string]interface{},
    ticket string,
    apiKey string,
) (map[string]interface{}, error) {
    
    cliente := NuevoClienteAlfresco(s.urlBase, apiKey, s.log)
    rawResponse, err := cliente.ListarDocumentos(c.Request.Context(), filtros, ticket)
    if err != nil {
        return nil, err
    }
    
    // Mapear a dominio
    documentosMap := make(map[string]interface{})
    for _, dto := range rawResponse {
        documentosMap[dto.Entry.ID] = documentos.Documento{
            ID: dto.Entry.ID,
            Nombre: dto.Entry.Name,
            Propiedades: dto.Entry.Properties,
        }
    }
    return documentosMap, nil
}

// DescargarDocumento implementa la descarga stateless.
// Parámetros:
//   - c: Contexto de Gin.
//   - idFile: ID del archivo a descargar.
//   - ticket: Ticket de autenticación.
//   - apiKey: API Key para autenticación.
// Retorna el contenido del archivo en bytes, el nombre del archivo y un error en caso de fallo.
func (s *ServicioAlfresco) DescargarDocumento(
    c *gin.Context,
    idFile string,
    ticket string,
    apiKey string,
) ([]byte, string, error) {

    cliente := NuevoClienteAlfresco(s.urlBase, apiKey, s.log)
    return cliente.DescargarDocumento(c.Request.Context(), idFile, ticket)
}