package servicio

import (
    "fmt"
    "mime/multipart"
    "github.com/gin-gonic/gin"
    "github.com/CamiloScript/REGAPIGO/shared/logger"
    "github.com/CamiloScript/REGAPIGO/application/documento"
)

// MockAuthClient simula el cliente de autenticación.
// Este mock permite probar el flujo de autenticación sin necesidad de conectarse a un servidor real de Alfresco.
type MockAuthClient struct {
    Log         *logger.Registrador // Logger para registrar eventos y errores
    ForzarError bool                // Bandera para forzar errores en las operaciones simuladas
}

// Authenticate simula el proceso de autenticación.
// Si ForzarError es true, devuelve un error simulado. De lo contrario, devuelve un ticket de autenticación mock.
func (m *MockAuthClient) Authenticate(usuario, password string) (string, error) {
    if m.ForzarError {
        m.Log.Error("MockAuthClient: Error forzado en Login", nil)
        return "", fmt.Errorf("error simulado")
    }
    m.Log.Info("MockAuthClient: Login simulado", map[string]interface{}{"usuario": usuario})
    return "TICKET_mock_123", nil
}

// MockClienteAlfresco simula todas las operaciones de documentos en Alfresco.
// Este mock permite probar el flujo de subida, listado y descarga de documentos sin necesidad de conectarse a un servidor real.
type MockClienteAlfresco struct {
    Log          *logger.Registrador               // Logger para registrar eventos y errores
    ForzarError  bool                              // Bandera para forzar errores en las operaciones simuladas
    DocumentosDB map[string]map[string]interface{} // Base de datos simulada de documentos
}

// NuevoMockClienteAlfresco crea una nueva instancia de MockClienteAlfresco.
// Inicializa una base de datos simulada con un documento de prueba.
func NuevoMockClienteAlfresco(log *logger.Registrador) *MockClienteAlfresco {
    return &MockClienteAlfresco{
        Log: log,
        DocumentosDB: map[string]map[string]interface{}{
            "doc-123": {
                "id":   "doc-123",
                "name": "contrato.pdf",
                "properties": map[string]interface{}{
                    "tanner:rut-cliente": "18079686-6",
                },
            },
        },
    }
}

// SubirDocumento simula la subida de un documento a Alfresco.
// Si ForzarError es true, devuelve un error simulado. De lo contrario, devuelve una respuesta mock.
func (m *MockClienteAlfresco) SubirDocumento(
    ctx *gin.Context,                  // Contexto de Gin
    archivo *multipart.FileHeader,     // Archivo a subir
    metadatos string,                  // Metadatos del documento
    ticket string,                     // Ticket de autenticación
    nombreArchivo string,              // Nombre del archivo
) (map[string]interface{}, error) {
    if m.ForzarError {
        m.Log.Error("MockClienteAlfresco: Error forzado en SubirDocumento", nil)
        return nil, fmt.Errorf("error simulado")
    }

    m.Log.Info("MockClienteAlfresco: Documento subido", map[string]interface{}{
        "filename": archivo.Filename,
        "ticket":   ticket,
    })

    // Retornar una respuesta mock con el ID y nombre del archivo subido
    return map[string]interface{}{
        "entry": map[string]interface{}{
            "id":   "mock-456",
            "name": archivo.Filename,
        },
    }, nil
}

// ListarDocumentos simula el listado de documentos en Alfresco.
// Si ForzarError es true, devuelve un error simulado. De lo contrario, filtra los documentos según los filtros proporcionados.
func (m *MockClienteAlfresco) ListarDocumentos(
    ctx *gin.Context,
    filtros map[string]interface{},
    ticket string,
    apiKey string,
) (map[string]interface{}, error) { 
    if m.ForzarError {
        m.Log.Error("MockClienteAlfresco: Error forzado en ListarDocumentos", nil)
        return nil, fmt.Errorf("error simulado")
    }

    // Filtrar documentos según el RUT del cliente
    resultados := []map[string]interface{}{}
    for _, doc := range m.DocumentosDB {
        if rut, ok := doc["properties"].(map[string]interface{})["tanner:rut-cliente"]; ok {
            if rut == filtros["tanner:rut-cliente"] {
                resultados = append(resultados, doc)
            }
        }
    }

    m.Log.Info("MockClienteAlfresco: Documentos listados", map[string]interface{}{
        "total": len(resultados),
    })

    // Retornar en el formato esperado por la interfaz
    return map[string]interface{}{
        "entries": resultados,
        "total": len(resultados),
    }, nil
}

// DescargarDocumento simula la descarga de un documento desde Alfresco.
// Si ForzarError es true, devuelve un error simulado. De lo contrario, devuelve el contenido mock del documento.
func (m *MockClienteAlfresco) DescargarDocumento(
    c *gin.Context,                    // Contexto de Gin
    idFile string,                     // ID del documento a descargar
    ticket string,                     // Ticket de autenticación
    otroParametro string,              // Parámetro adicional (no utilizado en este mock)
) ([]byte, string, error) {
    if m.ForzarError {
        m.Log.Error("MockClienteAlfresco: Error forzado en DescargarDocumento", nil)
        return nil, "", fmt.Errorf("error simulado")
    }

    // Buscar el documento en la base de datos simulada
    doc, existe := m.DocumentosDB[idFile]
    if !existe {
        m.Log.Warn("MockClienteAlfresco: Documento no encontrado", map[string]interface{}{"idFile": idFile})
        return nil, "", documento.ErrDocumentoNoEncontrado
    }

    m.Log.Info("MockClienteAlfresco: Documento descargado", map[string]interface{}{
        "idFile": idFile,
    })

    // Retornar contenido mock y el nombre del documento
    return []byte("contenido-mock"), doc["name"].(string), nil
}