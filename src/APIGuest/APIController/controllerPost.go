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
    "go.mongodb.org/mongo-driver/mongo"
)

// CreateGuest maneja la solicitud POST para crear un nuevo huésped
func CreateGuest(c *gin.Context) {
    // Log de inicio de la operación, para seguimiento y registro de actividad.
    log.Debug().
        Str("endpoint", c.Request.URL.Path).
        Str("method", c.Request.Method).
        Msg("Iniciando creación de huésped")

    // 1. Validación del paquete de datos Guest, dentro de struct, y manejo de errores
    var newGuest APIStruct.Guest
    // El método ShouldBindJSON se encarga de validar y mapear los datos JSON recibidos en la estructura Guest.
    if err := c.ShouldBindJSON(&newGuest); err != nil {
        // Log de error en formato de datos de entrada, (Datos inválidos o faltantes)
        log.Error().
            Err(err).
            Str("client_ip", c.ClientIP()).
            Interface("request_body", c.Request.Body).
            Msg("Error en formato de datos de entrada")
        // Respuesta de error al cliente, con mensaje personalizado.
        c.JSON(http.StatusBadRequest, FormatResponse(
            "error", 
            "Estructura de datos inválida. Verifique: " + err.Error(),
            nil))
        return
    }

    // 2. Log de paquete de datos recibido, para seguimiento y registro de actividad, (sin datos sensibles)
    log.Info().
        Str("email", newGuest.Email).
        Str("nombre", newGuest.FirstName + " " + newGuest.LastName).
        Msg("Payload recibido para creación de huésped")

    // 3. Asignación de metadatos temporales que estaran en la estructura Guest
    newGuest.CreatedAt = time.Now().UTC()
    newGuest.UpdatedAt = time.Now().UTC()
    
    // Log de asignación de metadatos temporales, para seguimiento y registro de actividad.
    log.Debug().
        Time("created_at", newGuest.CreatedAt).
        Time("updated_at", newGuest.UpdatedAt).
        Msg("Metadatos temporales asignados")

    // 4. Se genera la conexion con la base de datos y la colección de huéspedes, para la operación de inserción de info a la base de datos.
    collection := MongoDB.Cliente.Database("APIREGDB").Collection("Guests")
    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
    defer cancel()

    // 5. Se agregan los datos del nuevo huésped a la base de datos, y se manejan los errores de la operación.
    result, err := collection.InsertOne(ctx, newGuest)
    if err != nil {
        log.Error().
            Err(err).
            Interface("huésped", newGuest).
            Str("collection", "Guests").
            Msg("Fallo en operación MongoDB InsertOne")

        // Detección de errores específicos de MongoDB, para un mejor manejo y respuesta al cliente.
        if mongoErr, ok := err.(mongo.WriteException); ok {
            log.Warn().
                Int("code", mongoErr.WriteConcernError.Code).
                Interface("details", mongoErr.WriteErrors).
                Msg("Error específico de MongoDB")
        }
        // Respuesta de error al cliente, con mensaje personalizado, el cual se puede modificar en el archivo formatResponse.go
        c.JSON(http.StatusInternalServerError, FormatResponse(
            "error", 
            "Error interno al procesar la solicitud",
            nil))
        return
    }

    // 6. Se verifica el resultado de la operación de inserción, y se manejan los errores de la operación.
    if result.InsertedID == nil {
        log.Warn().Msg("InsertOne no devolvió un ID válido")
        c.JSON(http.StatusInternalServerError, FormatResponse(
            "error", 
            "Inconsistencia en el resultado de la operación",
            nil))
        return
    }

    // 7. Se genera un componente para realizar la construcción de la respuesta exitosa, eligiendo el tipo de dato a mostrar.
    responseData := map[string]interface{}{
        "guest_id": result.InsertedID,
        "email":    newGuest.Email,
    }
    // Log de creación exitosa de huésped, para seguimiento y registro de actividad.
    log.Info().
        Str("inserted_id", fmt.Sprintf("%v", result.InsertedID)).
        Dur("duración_operación", time.Since(newGuest.CreatedAt)).
        Msg("Huésped creado exitosamente")

    // 8. Se genera la respuesta exitosa al cliente, con mensaje personalizado, el cual se puede modificar en el archivo formatResponse.go
    c.JSON(http.StatusCreated, FormatResponse(
        "success", 
        "Huésped registrado exitosamente",
        responseData))
}
