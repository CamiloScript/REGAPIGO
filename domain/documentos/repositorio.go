package documentos

import (
    "github.com/gin-gonic/gin" // Framework web para manejar solicitudes HTTP
)

// AlmacenamientoDocumentos define las operaciones necesarias para interactuar con Alfresco.
// Esta interfaz es stateless, lo que significa que no almacena estado y requiere credenciales en cada solicitud.
type AlmacenamientoDocumentos interface {
    // SubirDocumento maneja la subida de un documento a Alfresco.
    // Recibe:
    // - c: Contexto de Gin para manejar la solicitud HTTP.
    // - fileBytes: Bytes del archivo a subir.
    // - metadatos: Metadatos asociados al archivo en formato string.
    // - ticket: Ticket de autenticación para Alfresco.
    // - apiKey: API Key de Alfresco para autorización.
    // Retorna:
    // - Un mapa con la respuesta de Alfresco.
    // - Un error en caso de fallo.
    SubirDocumento(
        c *gin.Context,
        fileBytes []byte,
        metadatos string,
        ticket string,
        apiKey string,
    ) (map[string]interface{}, error)


    // ListarDocumentos maneja la obtención de un listado de documentos desde Alfresco.
    // Recibe:
    // - c: Contexto de Gin para manejar la solicitud HTTP.
    // - filtros: Mapa de filtros para la búsqueda de documentos.
    // - ticket: Ticket de autenticación para Alfresco.
    // - nombreArchivo: Nombre del archivo para filtrar (opcional).
    // Retorna:
    // - Un mapa con la respuesta de Alfresco.
    // - Un error en caso de fallo.
    ListarDocumentos(
        c *gin.Context,
        filtros map[string]interface{},
        ticket string,
        nombreArchivo string,
    ) (map[string]interface{}, error)

    // DescargarDocumento maneja la descarga de un documento desde Alfresco.
    // Recibe:
    // - c: Contexto de Gin para manejar la solicitud HTTP.
    // - idFile: ID del archivo a descargar.
    // - ticket: Ticket de autenticación para Alfresco.
    // - apiKey: API Key de Alfresco para autorización.
    // Retorna:
    // - El contenido binario del archivo.
    // - El tipo MIME del archivo.
    // - Un error en caso de fallo.
    DescargarDocumento(
        c *gin.Context,
        idFile string,
        ticket string,
        apiKey string,
    ) ([]byte, string, error)
}