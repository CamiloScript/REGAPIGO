package handlers

import (
	"fmt"
	"net/http"
	"github.com/CamiloScript/REGAPIGO/infraestructure/persistence/db/mongo"
	"github.com/CamiloScript/REGAPIGO/shared/utils"
	"github.com/gin-gonic/gin"
)

// LoteDocumentos representa la estructura para recibir múltiples documentos en una sola solicitud.
// Esta estructura contiene una lista de documentos individuales y metadatos comunes que se aplican a todos los documentos.
type LoteDocumentos struct {
	Documentos   []SolicitudIndividual  `json:"documentos"`   // Lista de documentos individuales.
	DatosComunes map[string]interface{} `json:"datosComunes"` // Metadatos comunes para todos los documentos.
}

// SolicitudIndividual representa cada documento individual en el lote.
// Contiene el archivo del documento en base64 y los metadatos específicos para ese documento.
type SolicitudIndividual struct {
	Base64    string                 `json:"base64"`    // Archivo del documento en base64.
	Metadatos map[string]interface{} `json:"metadatos"` // Metadatos específicos del documento.
}

// ResultadoCarga representa el resultado de la carga de cada documento.
// Incluye el nombre del archivo, el estado de la carga (EXITOSO o ERROR) y un mensaje de error en caso de fallo.
type ResultadoCarga struct {
	NombreArchivo string `json:"nombreArchivo"`         // Nombre del archivo cargado.
	Estado        string `json:"estado"`                // Estado de la carga (EXITOSO o ERROR).
	Error         string `json:"error,omitempty"`       // Mensaje de error en caso de fallo.
	RazonSocial   string `json:"razonSocial,omitempty"` // Razón social del cliente (extraída de RazonSocialCliente).
	Rut           string `json:"rut,omitempty"`         // RUT del cliente (extraído de RUTCliente).
	IdArchivo     string `json:"idArchivo,omitempty"`   // ID del archivo subido en Alfresco.
	Base64        string `json:"base64,omitempty"`      // Archivo del documento en base64.
}

// ManejadorLoteDocumentos procesa la subida de un lote de documentos.
// @Summary Subir un lote de documentos
// @Description Procesa y sube múltiples documentos en una sola solicitud JSON.
// @Tags documentos
// @Accept json
// @Produce json
// @Param lote body LoteDocumentos true "Lote de documentos a subir"
// @Success 200 {object} gin.H "Respuesta con el resultado de la carga"
// @Failure 400 {object} gin.H "Error en el formato de la solicitud"
// @Failure 500 {object} gin.H "Error interno del servidor"
// @Router /subir-lote [post]
func (h *ManejadorDocumentos) ManejadorLoteDocumentos(c *gin.Context) {
    // 1. Autenticación interna
    ticket, err := h.internalAuth.AutenticarInternamente()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno de autenticación"})
        return
    }

    // 2. Parsear solicitud JSON
    var lote LoteDocumentos
    if err := c.ShouldBindJSON(&lote); err != nil {
        h.log.Error("Solicitud inválida", map[string]interface{}{"error": err.Error()})
        c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de solicitud incorrecto"})
        return
    }

    // 3. Inicializar slices para respuestas y errores
    resultados := make([]ResultadoCarga, 0, len(lote.Documentos))
    errores := make([]string, 0)

    // 4. Procesar cada documento en el lote
    for _, documento := range lote.Documentos {
        // Decodificar el archivo base64 a bytes
        fileBytes, err := utils.DecodeBase64(documento.Base64)
        if err != nil {
            errores = append(errores, fmt.Sprintf("Error al decodificar base64 para %s: %v", documento.Metadatos["tanner:nombre-doc"], err))
            continue
        }

        // Llamar al servicio para subir el documento
        respuesta, err := h.servicio.SubirDocumento(c, fileBytes, documento.Metadatos, ticket)
        if err != nil {
            errores = append(errores, fmt.Sprintf("Error al subir documento %s: %v", documento.Metadatos["tanner:nombre-doc"], err))
            continue
        }

        // Construir resultado de carga
        resultado := ResultadoCarga{
            NombreArchivo: documento.Metadatos["tanner:nombre-doc"].(string),
            Estado:        "EXITOSO",
            IdArchivo:     respuesta["entry"].(map[string]interface{})["id"].(string),
            RazonSocial:   documento.Metadatos["tanner:razon-social-cliente"].(string),
            Rut:           documento.Metadatos["tanner:rut-cliente"].(string),
            Base64:        documento.Base64,
        }

        // Agregar a resultados
        resultados = append(resultados, resultado)

        // Guardar en MongoDB
        if err := mongo.GuardarEnMongoDB(respuesta, h.log); err != nil {
            h.log.Error("Error en persistencia MongoDB", map[string]interface{}{"error": err.Error()})
        }
    }

    // 5. Construir respuesta para el cliente
    respuestaCliente := gin.H{
        "total_procesados": len(resultados),
        "documentos":       resultados,
    }

    if len(errores) > 0 {
        respuestaCliente["errores"] = errores
        respuestaCliente["total_errores"] = len(errores)
    }

    c.JSON(http.StatusOK, respuestaCliente)
    h.log.Info("Lote procesado", map[string]interface{}{
        "total_procesados": len(resultados),
        "total_errores":    len(errores),
    })
}
