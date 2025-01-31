package APIController

import (
    "fmt"
    "net/http"
    "time"
    "context"
    "github.com/CamiloScript/REGAPIGO/src/APIGuest/APIStruct"
    "github.com/CamiloScript/REGAPIGO/bd/MongoDB"
    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog/log"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

// CreateGuest maneja la solicitud POST para crear un nuevo huésped
func CreateGuest(c *gin.Context) {
    // Log de inicio de la operación con contexto
    log.Debug().
        Str("endpoint", c.Request.URL.Path).
        Str("method", c.Request.Method).
        Msg("Iniciando creación de huésped")

    // 1. Validación del Package Guest (Estructura de datos)
    var newGuest APIStruct.Guest
    if err := c.ShouldBindJSON(&newGuest); err != nil {
        log.Error().
            Err(err).
            Str("client_ip", c.ClientIP()).
            Interface("request_body", c.Request.Body).
            Msg("Error en formato de datos de entrada")
        
        c.JSON(http.StatusBadRequest, FormatResponse(
            "error", 
            "Estructura de datos inválida. Verifique: " + err.Error(),
            nil))
        return
    }

    // 2. Log del Package recibido (sin datos sensibles)
    log.Info().
        Str("email", newGuest.Email).
        Str("nombre", newGuest.FirstName + " " + newGuest.LastName).
        Msg("Payload recibido para creación de huésped")

    // 3. Asignación de metadatos temporales
    newGuest.CreatedAt = time.Now().UTC()
    newGuest.UpdatedAt = time.Now().UTC()
    
    log.Debug().
        Time("created_at", newGuest.CreatedAt).
        Time("updated_at", newGuest.UpdatedAt).
        Msg("Metadatos temporales asignados")

    // 4. Conexión a MongoDB
    collection := MongoDB.Cliente.Database("APIREGDB").Collection("Guests")
    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
    defer cancel()

    // 5. Operación de inserción
    result, err := collection.InsertOne(ctx, newGuest)
    if err != nil {
        log.Error().
            Err(err).
            Interface("huésped", newGuest).
            Str("collection", "Guests").
            Msg("Fallo en operación MongoDB InsertOne")

        // Detección específica de errores de MongoDB
        if mongoErr, ok := err.(mongo.WriteException); ok {
            log.Warn().
                Int("code", mongoErr.WriteConcernError.Code).
                Interface("details", mongoErr.WriteErrors).
                Msg("Error específico de MongoDB")
        }

        c.JSON(http.StatusInternalServerError, FormatResponse(
            "error", 
            "Error interno al procesar la solicitud",
            nil))
        return
    }

    // 6. Verificación del resultado
    if result.InsertedID == nil {
        log.Warn().Msg("InsertOne no devolvió un ID válido")
        c.JSON(http.StatusInternalServerError, FormatResponse(
            "error", 
            "Inconsistencia en el resultado de la operación",
            nil))
        return
    }

    // 7. Construcción de respuesta
    responseData := map[string]interface{}{
        "guest_id": result.InsertedID,
        "email":    newGuest.Email,
    }

    log.Info().
        Str("inserted_id", fmt.Sprintf("%v", result.InsertedID)).
        Dur("duración_operación", time.Since(newGuest.CreatedAt)).
        Msg("Huésped creado exitosamente")

    // 8. Respuesta exitosa
    c.JSON(http.StatusCreated, FormatResponse(
        "success", 
        "Huésped registrado exitosamente",
        responseData))
}
