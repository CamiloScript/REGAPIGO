package APIController

import (
    "net/http"
    "github.com/CamiloScript/REGAPIGO/bd/MongoDB"
    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog/log"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)


// DeleteGuest maneja la solicitud DELETE para eliminar un huésped por su ID
func DeleteGuest(c *gin.Context) {
    id := c.Param("id")
    log.Debug().Str("id", id).Msg("Iniciando eliminación de huésped")

    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        log.Error().Err(err).Str("id", id).Msg("ID inválido para eliminación")
        c.JSON(http.StatusBadRequest, FormatResponse("error", "ID inválido", nil))
        return
    }

    collection := MongoDB.Cliente.Database("APIREGDB").Collection("Guests")
    result, err := collection.DeleteOne(c, bson.M{"_id": objectID})
    if err != nil {
        log.Error().Err(err).Str("id", id).Msg("Error al eliminar de MongoDB")
        c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al eliminar el huésped", nil))
        return
    }

    if result.DeletedCount == 0 {
        log.Warn().Str("id", id).Msg("Huésped no encontrado para eliminación")
        c.JSON(http.StatusNotFound, FormatResponse("error", "Huésped no encontrado", nil))
        return
    }

    log.Info().Str("id", id).Msg("Huésped eliminado exitosamente")
    c.JSON(http.StatusOK, FormatResponse("success", "Huésped eliminado correctamente", nil))
}

