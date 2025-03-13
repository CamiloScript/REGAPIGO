package middleware

import (
    "time"
    "github.com/gin-gonic/gin"
    "github.com/CamiloScript/REGAPIGO/shared/logger"
    "github.com/google/uuid"
)

// MiddlewareRegistro crea un middleware de logging estructurado.
// Parámetros:
//   - registro: Instancia del registrador para registrar eventos y errores.
// Retorna una función de middleware para Gin.
func MiddlewareRegistro(registro *logger.Registrador) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Generar un ID único para la solicitud
        idSolicitud := uuid.New().String()
        c.Set("idSolicitud", idSolicitud)

        // Registrar el momento de inicio de la solicitud
        inicio := time.Now()

        // Continuar con la ejecución de los siguientes handlers en la cadena
        c.Next()

        // Calcular la duración total de la solicitud
        duracion := time.Since(inicio)

        // Construir un mapa con los datos de la solicitud y la respuesta para el registro
        contextoRegistro := map[string]interface{}{
            "id_solicitud":     idSolicitud,           // ID único de la solicitud
            "duracion":         duracion.String(),     // Duración total del procesamiento en formato de cadena
            "estado":           c.Writer.Status(),     // Código de estado HTTP devuelto en la respuesta
            "ip_cliente":       c.ClientIP(),          // Dirección IP del cliente que hizo la solicitud
            "metodo":           c.Request.Method,      // Método HTTP utilizado en la solicitud (GET, POST, etc.)
            "ruta":             c.Request.URL.Path,    // Ruta solicitada por el cliente
            "tamano_respuesta": c.Writer.Size(),       // Tamaño de la respuesta en bytes
        }

        // Si el encabezado User-Agent está presente, lo agrega al contexto del log
        if agenteUsuario := c.Request.Header.Get("User-Agent"); agenteUsuario != "" {
            contextoRegistro["agente_usuario"] = agenteUsuario
        }

        // Si el encabezado Referer está presente, lo agrega al contexto del log
        if referencia := c.Request.Header.Get("Referer"); referencia != "" {
            contextoRegistro["referencia"] = referencia
        }

        // Determinar el nivel de registro basado en el código de estado HTTP de la respuesta
        switch {
        case c.Writer.Status() >= 500:
            // Registra un error si el código de estado es 500 o superior (errores del servidor)
            registro.Error("Solicitud Finalizada, Error del Servidor", contextoRegistro)
        case c.Writer.Status() >= 400:
            // Registra una advertencia si el código de estado está entre 400 y 499 (errores del cliente)
            registro.Warn("Solicitud Finalizada, Error del Servicio", contextoRegistro)
        default:
            // Registra información estándar para respuestas exitosas (códigos 200-399)
            registro.Info("Solicitud Completada Exitosamente", contextoRegistro)
        }
    }
}