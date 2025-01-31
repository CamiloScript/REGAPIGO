package APIController

// Se importan las dependencias necesarias.
import (
    "net/http"
    "github.com/CamiloScript/REGAPIGO/bd/MongoDB"
    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog/log"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)


// DeleteGuest genera la eliminación de un huésped en la base de datos, mediante su ID, y retorna un mensaje de éxito o error.
func DeleteGuest(c *gin.Context) {
    // Se obtiene el ID del huésped a eliminar, y se registra con zerolog para su seguimiento.
    id := c.Param("id")
    log.Debug().Str("id", id).Msg("Iniciando eliminación de huésped")

    // Se transforma el ID del huésped a un ObjectID de MongoDB, y se registra con zerolog para su seguimiento.
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        log.Error().Err(err).Str("id", id).Msg("ID inválido para eliminación")
        c.JSON(http.StatusBadRequest, FormatResponse("error", "ID inválido", nil))
        return
    }
    // Se realiza una busqueda en la base de datos para verificar la existencia del huésped, y se registra con zerolog para su seguimiento.
    collection := MongoDB.Cliente.Database("APIREGDB").Collection("Guests")
    result, err := collection.DeleteOne(c, bson.M{"_id": objectID})
    // Se genera un registro con zerolog si ocurre un error al eliminar el huésped, y se retorna un mensaje de error.
    if err != nil {
        log.Error().Err(err).Str("id", id).Msg("Error al eliminar de MongoDB")
        c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al eliminar el huésped", nil))
        return
    }
    // Se genera un registro con zerolog si no se encuentra el huésped para eliminar, y se retorna un mensaje de error.
    if result.DeletedCount == 0 {
        log.Warn().Str("id", id).Msg("Huésped no encontrado para eliminación")
        c.JSON(http.StatusNotFound, FormatResponse("error", "Huésped no encontrado", nil))
        return
    }
    // Se genera un registro con zerolog si el huésped es eliminado exitosamente, y se retorna un mensaje de éxito.
    log.Info().Str("id", id).Msg("Huésped eliminado exitosamente")
    c.JSON(http.StatusOK, FormatResponse("success", "Huésped eliminado correctamente", nil))
}

