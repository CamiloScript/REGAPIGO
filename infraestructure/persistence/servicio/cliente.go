package servicio

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "regexp"
    "path/filepath"
    "time"
    "github.com/CamiloScript/REGAPIGO/shared/logger"
    "github.com/CamiloScript/REGAPIGO/application/documento"
)

// ClienteAlfresco maneja operaciones con Alfresco sin estado.
type ClienteAlfresco struct {
    urlBase     string              // URL base del servidor Alfresco
    clienteHTTP *http.Client        // Cliente HTTP con timeout
    log         *logger.Registrador // Logger para registro de eventos
    apiKey      string              // API Key para autenticación
}

// NuevoClienteAlfresco crea una instancia del cliente.
// Parámetros:
//   - urlBase: URL base del servidor Alfresco.
//   - apiKey: API Key para autenticación.
//   - log: Logger para registrar eventos y errores.
// Retorna una instancia configurada de ClienteAlfresco.
func NuevoClienteAlfresco(urlBase, apiKey string, log *logger.Registrador) *ClienteAlfresco {
    return &ClienteAlfresco{
        urlBase:     urlBase,
        apiKey:      apiKey,
        clienteHTTP: &http.Client{Timeout: 30 * time.Second},
        log:         log,
    }
}

// SubirDocumento envía un archivo y metadatos a Alfresco.
// Parámetros:
//   - ctx: Contexto para controlar la solicitud.
//   - archivo: Archivo a subir (multipart).
//   - metadatos: Metadatos del archivo en formato JSON.
//   - ticket: Ticket de autenticación.
// Retorna un mapa con la respuesta de Alfresco o un error en caso de fallo.
func (c *ClienteAlfresco) SubirDocumento(
    ctx context.Context,
    archivo *multipart.FileHeader,
    metadatos string,
    ticket string,
) (map[string]interface{}, error) {

    endpoint := "/tanner-alfresco/file-upload"
    url := c.urlBase + endpoint

    // Preparar formulario multipart
    cuerpo := &bytes.Buffer{}
    escritor := multipart.NewWriter(cuerpo)

    // Abrir archivo
    f, err := archivo.Open()
    if err != nil {
        c.log.Error("Error al abrir archivo", map[string]interface{}{"error": err.Error()})
        return nil, fmt.Errorf("error al abrir archivo: %v", err)
    }
    defer f.Close()

    // Crear parte del archivo
    parte, err := escritor.CreateFormFile("documento", filepath.Base(archivo.Filename))
    if err != nil {
        c.log.Error("Error al crear formulario", map[string]interface{}{"error": err.Error()})
        return nil, fmt.Errorf("error al crear formulario: %v", err)
    }

    // Copiar contenido del archivo
    if _, err := io.Copy(parte, f); err != nil {
        c.log.Error("Error al copiar archivo", map[string]interface{}{"error": err.Error()})
        return nil, fmt.Errorf("error al copiar archivo: %v", err)
    }

    // Agregar metadatos
    if err := escritor.WriteField("propiedades", metadatos); err != nil {
        c.log.Error("Error al escribir metadatos", map[string]interface{}{"error": err.Error()})
        return nil, fmt.Errorf("error al escribir metadatos: %v", err)
    }
    escritor.Close()

    // Crear solicitud HTTP
    req, err := http.NewRequestWithContext(ctx, "POST", url, cuerpo)
    if err != nil {
        c.log.Error("Error al crear solicitud", map[string]interface{}{"error": err.Error()})
        return nil, fmt.Errorf("error al crear solicitud: %v", err)
    }

    // Configurar headers
    req.Header.Set("Content-Type", escritor.FormDataContentType())
    req.Header.Set("Authorization", "Basic "+ticket) // Ticket dinámico
    req.Header.Set("ADFTannerServices", c.apiKey)    // API Key fija

    // Enviar solicitud
    var resultado map[string]interface{}
    if err := c.ejecutarSolicitud(req, &resultado); err != nil {
        c.log.Error("Error al subir archivo", map[string]interface{}{"error": err.Error()})
        return nil, fmt.Errorf("error al subir archivo: %v", err)
    }

    return resultado, nil
}

// ListarDocumentos recupera documentos filtrados desde Alfresco.
// Parámetros:
//   - ctx: Contexto para controlar la solicitud.
//   - filtros: Filtros de búsqueda en formato JSON.
//   - ticket: Ticket de autenticación.
// Retorna una lista de documentos de Alfresco o un error en caso de fallo.
func (c *ClienteAlfresco) ListarDocumentos(
    ctx context.Context,
    filtros map[string]interface{},
    ticket string,
) ([]documento.AlfrescoDocumentDTO, error) {
    endpoint := "/tanner-alfresco/files"
    url := c.urlBase + endpoint

    // Convertir filtros a JSON
    cuerpo, err := json.Marshal(filtros)
    if err != nil {
        c.log.Error("Error al serializar filtros", map[string]interface{}{"error": err.Error()})
        return nil, fmt.Errorf("error al serializar filtros: %v", err)
    }

    // Crear solicitud HTTP
    req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(cuerpo))
    if err != nil {
        c.log.Error("Error al crear solicitud", map[string]interface{}{"error": err.Error()})
        return nil, fmt.Errorf("error al crear solicitud: %v", err)
    }

    // Configurar headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Basic "+ticket) // Ticket dinámico
    req.Header.Set("ADFTannerServices", c.apiKey)    // API Key fija

    // Ejecutar solicitud
    var respuesta []documento.AlfrescoDocumentDTO
    if err := c.ejecutarSolicitud(req, &respuesta); err != nil {
        return nil, err
    }
    return respuesta, nil
}

// DescargarDocumento obtiene un archivo por su ID.
// Parámetros:
//   - ctx: Contexto para controlar la solicitud.
//   - idFile: ID del archivo a descargar.
//   - ticket: Ticket de autenticación.
// Retorna el contenido del archivo en bytes, el nombre del archivo y un error en caso de fallo.
func (c *ClienteAlfresco) DescargarDocumento(
    ctx context.Context,
    idFile string,
    ticket string,
) ([]byte, string, error) {
    endpoint := "/tanner-alfresco/file-download?idFile=" + idFile
    url := c.urlBase + endpoint

    // Crear solicitud HTTP
    req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
    if err != nil {
        c.log.Error("Error al crear solicitud", map[string]interface{}{"error": err.Error()})
        return nil, "", fmt.Errorf("error al crear solicitud: %v", err)
    }

    // Configurar headers
    req.Header.Set("Authorization", "Basic "+ticket)
    req.Header.Set("ADFTannerServices", c.apiKey)

    // Ejecutar solicitud
    resp, err := c.clienteHTTP.Do(req)
    if err != nil {
        c.log.Error("Error al descargar archivo", map[string]interface{}{"error": err.Error()})
        return nil, "", fmt.Errorf("error al descargar archivo: %v", err)
    }
    defer resp.Body.Close()

    // Manejar errores HTTP
    if resp.StatusCode >= 400 {
        cuerpoError, _ := io.ReadAll(resp.Body)
        c.log.Error("Error de Alfresco", map[string]interface{}{
            "status_code": resp.StatusCode,
            "respuesta":   string(cuerpoError),
        })
        return nil, "", fmt.Errorf("error de Alfresco: %s. Detalle: %s", resp.Status, cuerpoError)
    }

    // Leer contenido del archivo
    contenido, err := io.ReadAll(resp.Body)
    if err != nil {
        c.log.Error("Error al leer contenido", map[string]interface{}{"error": err.Error()})
        return nil, "", fmt.Errorf("error al leer contenido: %v", err)
    }

    // Extraer nombre del archivo del header
    nombreArchivo := "documento_" + idFile
    if contentDisposition := resp.Header.Get("Content-Disposition"); contentDisposition != "" {
        re := regexp.MustCompile(`filename="(.+?)"`)
        matches := re.FindStringSubmatch(contentDisposition)
        if len(matches) > 1 {
            nombreArchivo = matches[1]
        }
    }

    return contenido, nombreArchivo, nil
}

// ejecutarSolicitud maneja la lógica común para enviar solicitudes HTTP y procesar respuestas.
// Parámetros:
//   - req: Solicitud HTTP a ejecutar.
//   - destino: Estructura donde se decodificará la respuesta JSON.
// Retorna un error en caso de fallo.
func (c *ClienteAlfresco) ejecutarSolicitud(req *http.Request, destino interface{}) error {
    // Ejecutar la solicitud HTTP
    resp, err := c.clienteHTTP.Do(req)
    if err != nil {
        // Registrar error en el log
        c.log.Error("Error en la solicitud HTTP", map[string]interface{}{"error": err.Error()})
        return fmt.Errorf("error en la solicitud HTTP: %v", err)
    }
    defer resp.Body.Close() // Asegurar que el cuerpo de la respuesta se cierre

    // Verificar si la respuesta es un error (código 4xx o 5xx)
    if resp.StatusCode >= 400 {
        // Leer el cuerpo del error para obtener detalles
        cuerpoError, err := io.ReadAll(resp.Body)
        if err != nil {
            c.log.Error("Error al leer cuerpo de error", map[string]interface{}{"error": err.Error()})
            return fmt.Errorf("error de Alfresco: %s", resp.Status)
        }

        // Registrar el error en el log
        c.log.Error("Error de Alfresco", map[string]interface{}{
            "status_code": resp.StatusCode,
            "respuesta":   string(cuerpoError),
        })

        // Devolver un error con detalles
        return fmt.Errorf("error de Alfresco: %s. Detalle: %s", resp.Status, cuerpoError)
    }

    // Decodificar la respuesta JSON en la estructura destino
    if err := json.NewDecoder(resp.Body).Decode(destino); err != nil {
        // Registrar error en el log
        c.log.Error("Error al decodificar respuesta", map[string]interface{}{"error": err.Error()})
        return fmt.Errorf("error al decodificar respuesta: %v", err)
    }

    // Si todo está bien, retornar nil (sin errores)
    return nil
}