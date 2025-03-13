
package handlers

import (
    "net/http"
    "github.com/CamiloScript/REGAPIGO/application/documento"
    "github.com/CamiloScript/REGAPIGO/infraestructure/persistence/db/mongo"
    "github.com/CamiloScript/REGAPIGO/shared/logger"
    "github.com/CamiloScript/REGAPIGO/shared/config"
    "github.com/gin-gonic/gin"
    "github.com/CamiloScript/REGAPIGO/domain/auth"
    "github.com/CamiloScript/REGAPIGO/shared/utils"
)

// ==================================================
// Estructuras y Constructor
// ==================================================

// ManejadorBusquedaDescarga gestiona operaciones de búsqueda y descarga.
type ManejadorBusquedaDescarga struct {
    servicio     *documento.ImplementacionServicioDocumentos // Servicio de documentos
    log          *logger.Registrador // Logger para registro de eventos
    cfg          *config.Config // Configuración de la aplicación
    internalAuth *InternalAuth // Autenticación interna
}

// NuevoManejadorBusquedaDescarga crea una instancia del manejador con dependencias inyectadas.
func NuevoManejadorBusquedaDescarga(
    servicio *documento.ImplementacionServicioDocumentos,
    log *logger.Registrador,
    cfg *config.Config,
    authServicio auth.AuthService, // Servicio de autenticación
) *ManejadorBusquedaDescarga {
    return &ManejadorBusquedaDescarga{
        servicio:     servicio,
        log:          log,
        cfg:          cfg,
        internalAuth: NewInternalAuth(authServicio, log, cfg), // Inicializar autenticación interna
    }
}

// ==================================================
// Definición de Solicitud
// ==================================================

// SolicitudBusqueda define el formato esperado para las solicitudes de búsqueda
type SolicitudBusqueda struct {
    // Formato guardado en MongoDB
    TipoDocumento   string `json:"tipo_documento,omitempty"`   // Campo opcional
    RUTCliente      string `json:"rut_cliente,omitempty"`      // Campo obligatorio
    NombreDocumento string `json:"nombre_documento,omitempty"` // Campo opcional
    FechaCarga      string `json:"fecha_carga,omitempty"`      // Campo opcional
    RazonSocial     string `json:"razon_social,omitempty"`     // Campo opcional

    // Formato del servicio externo
    TannerNombreDoc       string `json:"tanner:nombre-doc,omitempty"`
    TannerTipoDocumento   string `json:"tanner:tipo-documento,omitempty"`
    TannerRazonSocial     string `json:"tanner:razon-social-cliente,omitempty"`
    TannerRUTCliente      string `json:"tanner:rut-cliente,omitempty"`
    TannerEstadoVigencia  string `json:"tanner:estado-vigencia,omitempty"`
    TannerFechaCarga      string `json:"tanner:fecha-carga,omitempty"`
}


// construirFiltro construye un filtro de búsqueda para MongoDB basado en la solicitud
func construirFiltro(solicitud SolicitudBusqueda) map[string]interface{} {
    filtro := make(map[string]interface{})

    // Campos obligatorios
    if solicitud.RUTCliente != "" {
        filtro["metadatos.rut_cliente"] = solicitud.RUTCliente
    } else if solicitud.TannerRUTCliente != "" {
        filtro["metadatos.rut_cliente"] = solicitud.TannerRUTCliente
    }

    // Campos opcionales
    if solicitud.TipoDocumento != "" {
        filtro["metadatos.tipo_documento"] = solicitud.TipoDocumento
    } else if solicitud.TannerTipoDocumento != "" {
        filtro["metadatos.tipo_documento"] = solicitud.TannerTipoDocumento
    }

    if solicitud.NombreDocumento != "" {
        filtro["metadatos.nombre_documento"] = solicitud.NombreDocumento
    } else if solicitud.TannerNombreDoc != "" {
        filtro["metadatos.nombre_documento"] = solicitud.TannerNombreDoc
    }

    if solicitud.FechaCarga != "" {
        filtro["metadatos.fecha_carga"] = solicitud.FechaCarga
    } else if solicitud.TannerFechaCarga != "" {
        filtro["metadatos.fecha_carga"] = solicitud.TannerFechaCarga
    }

    if solicitud.RazonSocial != "" {
        filtro["metadatos.razon_social_cliente"] = solicitud.RazonSocial
    } else if solicitud.TannerRazonSocial != "" {
        filtro["metadatos.razon_social_cliente"] = solicitud.TannerRazonSocial
    }

    if solicitud.TannerEstadoVigencia != "" {
        filtro["metadatos.estado_vigencia"] = solicitud.TannerEstadoVigencia
    }

    return filtro
}



// ==================================================
// Manejador Principal
// ==================================================

// BuscarYDescargarDocumento maneja el flujo completo de búsqueda y descarga
// @Summary Busca un documento en MongoDB y lo descarga desde Alfresco
// @Description Recibe un RUT, busca en MongoDB, y descarga el archivo asociado desde Alfresco
// @Tags Documentos
// @Accept json
// @Produce octet-stream
// @Param solicitud body SolicitudBusqueda true "Criterios de búsqueda"
// @Security ApiKeyAuth
// @Success 200 {file} binary "Archivo descargado"
// @Failure 400 {object} map[string]string "Error en formato de solicitud"
// @Failure 404 {object} map[string]string "Documento no encontrado"
// @Failure 500 {object} map[string]string "Error interno"
// @Router /documentos/buscar-descargar [post]
// busqueda_descarga_handler.go

func (h *ManejadorBusquedaDescarga) BuscarYDescargarDocumento(c *gin.Context) {
    // 1. Autenticación interna
    ticket, err := h.internalAuth.AutenticarInternamente()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno de autenticación"})
        return
    }

    // 2. Validar y parsear solicitud
    var solicitud SolicitudBusqueda
    if err := c.ShouldBindJSON(&solicitud); err != nil {
        h.log.Error("Solicitud inválida", map[string]interface{}{"error": err.Error()})
        c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de solicitud incorrecto"})
        return
    }

    // 3. Construir filtro y buscar en MongoDB
    filtro := construirFiltro(solicitud)
    idFile, err := mongo.BuscarDocumento(filtro)
    if err != nil {
        h.log.Error("Documento no encontrado en MongoDB", map[string]interface{}{"filtro": filtro, "error": err.Error()})
        c.JSON(http.StatusNotFound, gin.H{"error": "No se encontraron documentos con los criterios proporcionados"})
        return
    }

    // 4. Descargar desde Alfresco
    contenido, nombreArchivo, err := h.servicio.DescargarDocumento(c, idFile, ticket)
    if err != nil {
        h.log.Error("Fallo en descarga desde Alfresco", map[string]interface{}{"idFile": idFile, "error": err.Error()})
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recuperar el archivo desde el repositorio"})
        return
    }

    // 5. Codificar el archivo a base64
    base64File := utils.EncodeToBase64(contenido)

    // 6. Configurar respuesta
    c.JSON(http.StatusOK, gin.H{
        "fileName": nombreArchivo,
        "base64":   base64File,
    })
    h.log.Info("Descarga exitosa", map[string]interface{}{"idFile": idFile, "nombreArchivo": nombreArchivo})
}