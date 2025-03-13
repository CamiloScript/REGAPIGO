package handlers

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
    "fmt"
    "io"
	"github.com/CamiloScript/REGAPIGO/domain/documentos"
	"github.com/gin-gonic/gin"
    "github.com/CamiloScript/REGAPIGO/infraestructure/persistence/db/mongo"
    "github.com/CamiloScript/REGAPIGO/shared/utils"
)

// LoteDocumentos representa la estructura para recibir múltiples documentos en una sola solicitud.
// Esta estructura contiene una lista de documentos individuales y metadatos comunes que se aplican a todos los documentos.
type LoteDocumentos struct {
    Documentos   []SolicitudIndividual      `json:"documentos"`   // Lista de documentos individuales a cargar.
    DatosComunes documentos.DocumentMetadata `json:"datosComunes"` // Metadatos comunes para todos los documentos.
}

// SolicitudIndividual representa cada documento individual en el lote.
// Contiene el archivo del documento y los metadatos específicos para ese documento.
type SolicitudIndividual struct {
    Base64    string                     `json:"base64"`    // Base64 del documento.
    Archivo   *multipart.FileHeader      `json:"archivo"`   // Archivo del documento.
    Metadatos documentos.DocumentMetadata `json:"metadatos"` // Metadatos específicos del documento.
}

// ResultadoCarga representa el resultado de la carga de cada documento.
// Incluye el nombre del archivo, el estado de la carga (EXITOSO o ERROR) y un mensaje de error en caso de fallo.
type ResultadoCarga struct {
    NombreArchivo string `json:"nombreArchivo"` // Nombre del archivo cargado.
    Estado        string `json:"estado"`        // Estado de la carga (EXITOSO o ERROR).
    Error         string `json:"error,omitempty"` // Mensaje de error en caso de fallo.
    RazonSocial   string `json:"razonSocial,omitempty"` // Razón social del cliente (extraída de RazonSocialCliente).
    Rut           string `json:"rut,omitempty"`         // RUT del cliente (extraído de RUTCliente).
    IdArchivo     string `json:"idArchivo,omitempty"`   // ID del archivo subido en Alfresco.
}

// ManejadorLoteDocumentos procesa la subida de un lote de documentos.
// Este método maneja la solicitud POST para subir múltiples documentos al sistema.
// ManejadorLoteDocumentos procesa la subida de un lote de documentos.
func (h *ManejadorDocumentos) ManejadorLoteDocumentos(c *gin.Context) {
    // 1. Autenticación interna
    ticket, err := h.internalAuth.AutenticarInternamente()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno de autenticación"})
        return
    }

    // 2. Procesar formulario con múltiples archivos
    form, err := c.MultipartForm()
    if err != nil {
        h.log.Error("Error al procesar formulario", map[string]interface{}{"error": err.Error()})
        c.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar el formulario"})
        return
    }

    // 3. Obtener archivos y metadatos
    archivos := form.File["documento"]
    if len(archivos) == 0 {
        h.log.Error("No se encontraron archivos", nil)
        c.JSON(http.StatusBadRequest, gin.H{"error": "No se encontraron archivos"})
        return
    }

    // 4. Obtener metadatos para cada archivo
    metadatosArray := form.Value["propiedades"]
    if len(metadatosArray) == 0 {
        h.log.Error("No se encontraron metadatos", nil)
        c.JSON(http.StatusBadRequest, gin.H{"error": "No se encontraron metadatos"})
        return
    }

    // 5. Parsear metadatos
    var metadatosLote []map[string]interface{}
    if err := json.Unmarshal([]byte(metadatosArray[0]), &metadatosLote); err != nil {
        h.log.Error("Error al parsear metadatos", map[string]interface{}{"error": err.Error()})
        c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de metadatos inválido"})
        return
    }

    // 6. Validar que la cantidad de metadatos coincide con la cantidad de archivos
    if len(metadatosLote) != len(archivos) {
        h.log.Error("Cantidad de metadatos no coincide con archivos", 
            map[string]interface{}{
                "archivos": len(archivos), 
                "metadatos": len(metadatosLote),
            })
        c.JSON(http.StatusBadRequest, gin.H{"error": "La cantidad de metadatos debe coincidir con la cantidad de archivos"})
        return
    }

    // 7. Inicializar slices para respuestas y errores
    resultados := make([]gin.H, 0, len(archivos))
    errores := make([]string, 0)

    // 8. Procesar cada archivo con sus metadatos
    for i, archivo := range archivos {
        // Convertir metadatos a JSON para cada archivo
        metadatosJSON, err := json.Marshal(metadatosLote[i])
        if err != nil {
            errores = append(errores, fmt.Sprintf("Error al procesar metadatos para %s: %v", archivo.Filename, err))
            continue
        }

        // Llamar al servicio para subir el documento
        var metadatosMap map[string]interface{}
        if err := json.Unmarshal(metadatosJSON, &metadatosMap); err != nil {
            errores = append(errores, fmt.Sprintf("Error al parsear metadatos para %s: %v", archivo.Filename, err))
            continue
        }
        file, err := archivo.Open()
        if err != nil {
            errores = append(errores, fmt.Sprintf("Error al abrir archivo %s: %v", archivo.Filename, err))
            continue
        }
        defer file.Close()

        fileContent, err := io.ReadAll(file)
        if err != nil {
            errores = append(errores, fmt.Sprintf("Error al leer archivo %s: %v", archivo.Filename, err))
            continue
        }

        respuesta, err := h.servicio.SubirDocumento(c, fileContent, metadatosMap, ticket)
        if err != nil {
            errores = append(errores, fmt.Sprintf("Error al subir archivo %s: %v", archivo.Filename, err))
            continue
        }

        // Codificar el archivo a base64
        base64File := utils.EncodeToBase64(fileContent)

        // Agregar a resultados
        resultados = append(resultados, gin.H{
            "nombre_archivo": archivo.Filename,
            "respuesta": gin.H{
                "fileName": archivo.Filename,
                "base64":   base64File,
            },
        })

        if err := mongo.GuardarEnMongoDB(respuesta, h.log); err != nil {
            h.log.Error("Error en persistencia MongoDB", map[string]interface{}{"error": err.Error()})
        }
    }

    // 9. Construir respuesta para el cliente
    respuestaCliente := gin.H{
        "total_procesados": len(resultados),
        "documentos": resultados,
    }
    
    if len(errores) > 0 {
        respuestaCliente["errores"] = errores
        respuestaCliente["total_errores"] = len(errores)
    }

    c.JSON(http.StatusOK, respuestaCliente)
    h.log.Info("Lote procesado", map[string]interface{}{
        "total_procesados": len(resultados),
        "total_errores": len(errores),
    })
}