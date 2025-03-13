package test_handlers

import (
    "bytes"
    "encoding/json"
    "io"
    "mime/multipart"
    "net/http"
    "net/http/httptest"
    "os"
    "path/filepath"
    "testing"
    "github.com/CamiloScript/REGAPIGO/application/documento"
    "github.com/CamiloScript/REGAPIGO/domain/documentos"
    "github.com/CamiloScript/REGAPIGO/infraestructure/api/handlers"
    "github.com/CamiloScript/REGAPIGO/infraestructure/persistence/servicio"
    "github.com/CamiloScript/REGAPIGO/shared/logger"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/CamiloScript/REGAPIGO/shared/config"
)

// TestSubirDocumentoExitoso - Prueba el caso exitoso de subida de un documento.
// Este test verifica que el endpoint de subida de documentos funcione correctamente cuando se proporcionan un archivo y metadatos válidos.
func TestSubirDocumentoExitoso(t *testing.T) {
    // Inicializar el logger para las pruebas
    log := logger.NuevoRegistrador("TEST", "|")

    // Crear un mock del cliente de Alfresco
    mockCliente := servicio.NuevoMockClienteAlfresco(log)

    // Crear el servicio de documentos utilizando el mock
    servicio := documento.NuevoServicioDocumentos(mockCliente, log, "mock-key")

    // Crear la configuración mock para las pruebas
    config := &config.Config{}

    // Crear un mock del servicio de autenticación
    // Nota: Si MockAuthClient no está definido en el paquete servicio, debes definirlo o usar una implementación real.
    mockAuth := &servicio.MockAuthClient{Log: log}

    // Crear el manejador de documentos utilizando el servicio y el mock de autenticación
    manejador := handlers.NuevoManejadorDocumentos(servicio, log, config, mockAuth)

    // Configurar el router de Gin para la prueba
    router := gin.Default()
    router.POST("/subir", manejador.ManejadorSubirDocumento)

    // Crear un buffer para almacenar el cuerpo de la solicitud multipart
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    log.Debug("Formulario multipart creado", map[string]interface{}{"boundary": writer.Boundary()})

    // Agregar archivo al formulario
    filePath := filepath.Join("testdata", "sample1.pdf")
    file, err := os.Open(filePath)
    if err != nil {
        t.Fatalf("Error al abrir el archivo de prueba: %v", err)
    }
    defer file.Close()

    part, err := writer.CreateFormFile("documento", "sample1.pdf")
    if err != nil {
        t.Fatalf("Error al crear el campo de archivo en el formulario: %v", err)
    }
    if _, err := io.Copy(part, file); err != nil {
        t.Fatalf("Error al copiar el contenido del archivo: %v", err)
    }
    log.Info("Archivo agregado al formulario", map[string]interface{}{"ruta": filePath})

    // Agregar metadatos al formulario
    metadatos := documentos.DocumentMetadata{
        CmTitle:        "Test",
        RUTCliente:    "20218874-5",
        TipoDocumento:  "Factura",
        EstadoVigencia: "Vigente",
        NombreDoc:      "sample1.pdf",
        CmVersionType:  "MAJOR",
        CmVersionLabel: "1.0",
        CmDescription:  "Documento de prueba",
        SubCategorias:  "Riesgo", // Añadir sub-categoría válida
    }
    jsonMeta, err := json.Marshal(metadatos)
    if err != nil {
        t.Fatalf("Error al serializar los metadatos: %v", err)
    }
    if err := writer.WriteField("propiedades", string(jsonMeta)); err != nil {
        t.Fatalf("Error al agregar los metadatos al formulario: %v", err)
    }
    writer.Close()

    // Crear la solicitud HTTP POST
    w := httptest.NewRecorder()
    req, err := http.NewRequest("POST", "/subir", body)
    if err != nil {
        t.Fatalf("Error al crear la solicitud HTTP: %v", err)
    }
    req.Header.Set("Content-Type", writer.FormDataContentType())
    req.Header.Set("Authorization", "Basic TICKET_mock_123")

    // Ejecutar la solicitud
    router.ServeHTTP(w, req)

    // Verificar que el código de estado HTTP sea 200 (OK)
    assert.Equal(t, http.StatusOK, w.Code)
}

// TestListarDocumentos_DeAlfresco - Prueba el caso exitoso de listar documentos desde Alfresco.
// Este test verifica que el método ListarDocumentos funcione correctamente cuando se proporcionan filtros válidos.
func TestListarDocumentos_DeAlfresco(t *testing.T) {
    // Inicializar el logger
    log := logger.NuevoRegistrador("TEST", "|")
    
    // Crear un mock del cliente de Alfresco
    repo := servicio.NuevoMockClienteAlfresco(log)
    repo.ForzarError = false
    
    // Crear un mock del servicio de autenticación
    // Nota: Si MockAuthClient no está definido en el paquete servicio, debes definirlo o usar una implementación real.
    mockAuth := &servicio.MockAuthClient{Log: log}
    
    // Crear el servicio de documentos utilizando el mock
    servicio := documento.NuevoServicioDocumentos(repo, log, "mock-key")
    
    // Crear la configuración mock para las pruebas
    config := &config.Config{}
    
    // Crear el manejador de documentos utilizando el servicio y el mock de autenticación
    manejador := handlers.NuevoManejadorDocumentos(servicio, log, config, mockAuth)
    
    // Configurar el router de Gin para la prueba
    router := gin.Default()
    router.POST("/listar", manejador.ManejadorListarDocumentos)
    
    // Crear la solicitud HTTP POST con un cuerpo JSON que contiene filtros
    w := httptest.NewRecorder()
    reqBody := bytes.NewBufferString(`{"tanner:rut-cliente": "18079686-6"}`)
    req, err := http.NewRequest("POST", "/listar", reqBody)
    if err != nil {
        t.Fatalf("Error al crear la solicitud HTTP: %v", err)
    }
    req.Header.Set("Authorization", "Basic TICKET_mock_123")
    req.Header.Set("Content-Type", "application/json")
    
    // Ejecutar la solicitud
    router.ServeHTTP(w, req)
    
    // Verificar que el código de estado HTTP sea 200 (OK)
    assert.Equal(t, http.StatusOK, w.Code)
    
    // Verificar que la respuesta contenga el documento esperado
    assert.Contains(t, w.Body.String(), "doc-123")
}

// TestDescargarDocumentoNoExistente - Prueba el caso de error cuando se intenta descargar un documento que no existe.
// Este test verifica que el endpoint de descarga devuelva un error 404 (Not Found) cuando el documento no existe.
func TestDescargarDocumentoNoExistente(t *testing.T) {
    // Inicializar el logger para las pruebas
    log := logger.NuevoRegistrador("TEST", "|")
    
    // Crear un mock del cliente de Alfresco
    mockCliente := servicio.NuevoMockClienteAlfresco(log)
    
    // Crear un mock del servicio de autenticación
    // Nota: Si MockAuthClient no está definido en el paquete servicio, debes definirlo o usar una implementación real.
    mockAuth := &servicio.MockAuthClient{Log: log}
    
    // Crear la configuración mock para las pruebas
    config := &config.Config{}
    
    // Crear el servicio de documentos utilizando el mock
    servicio := documento.NuevoServicioDocumentos(mockCliente, log, "mock-key")
    
    // Crear el manejador de documentos utilizando el servicio y el mock de autenticación
    manejador := handlers.NuevoManejadorDocumentos(servicio, log, config, mockAuth)
    
    // Configurar el router de Gin para la prueba
    router := gin.Default()
    router.POST("/descargar", manejador.ManejadorDescargarDocumento)
    
    // Crear la solicitud HTTP POST para descargar un documento que no existe
    req, err := http.NewRequest("POST", "/descargar?idFile=no-existe", nil)
    if err != nil {
        t.Fatalf("Error al crear la solicitud HTTP: %v", err)
    }
    req.Header.Set("Authorization", "Basic TICKET_mock_123")
    
    // Ejecutar la solicitud
    resp := httptest.NewRecorder()
    log.Info("Ejecutando solicitud HTTP", nil)
    router.ServeHTTP(resp, req)
    
    log.Info("Respuesta del servidor", map[string]interface{}{
        "status":  resp.Code,
        "headers": resp.Header(),
    })
    
    // Verificar que el código de estado HTTP sea 404 (Not Found)
    assert.Equal(t, http.StatusNotFound, resp.Code)
    
    // Verificar que el tipo de contenido de la respuesta sea JSON
    assert.Equal(t, "application/json; charset=utf-8", resp.Header().Get("Content-Type"))
    
    // Verificar que la respuesta contenga el mensaje de error esperado
    assert.Contains(t, resp.Body.String(), "documento no encontrado")
    log.Info("Afirmaciones completadas", nil)
}