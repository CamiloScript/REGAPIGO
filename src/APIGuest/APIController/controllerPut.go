package APIController

import (
    "net/http"
    "time"
    "github.com/CamiloScript/REGAPIGO/src/APIGuest/APIStruct"
    "github.com/CamiloScript/REGAPIGO/bd/MongoDB"
    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog/log"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateGuest maneja la solicitud PUT para actualizar los datos de un huésped
func UpdateGuest(c *gin.Context) {
    id := c.Param("id")
    log.Debug().Str("id", id).Msg("Iniciando actualización de huésped")

    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        log.Error().Err(err).Str("id", id).Msg("ID inválido para actualización")
        c.JSON(http.StatusBadRequest, FormatResponse("error", "ID inválido", nil))
        return
    }

    var updatedGuest APIStruct.Guest
    if err := c.ShouldBindJSON(&updatedGuest); err != nil {
        log.Error().Err(err).Msg("Datos inválidos en solicitud de actualización")
        c.JSON(http.StatusBadRequest, FormatResponse("error", "Datos inválidos", nil))
        return
    }

    updatedGuest.UpdatedAt = time.Now()

    collection := MongoDB.Cliente.Database("APIREGDB").Collection("Guests")
    result, err := collection.UpdateOne(c, bson.M{"_id": objectID}, bson.M{"$set": updatedGuest})
    if err != nil {
        log.Error().Err(err).Str("id", id).Msg("Error al actualizar en MongoDB")
        c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al actualizar el huésped", nil))
        return
    }

    if result.MatchedCount == 0 {
        log.Warn().Str("id", id).Msg("Huésped no encontrado para actualización")
        c.JSON(http.StatusNotFound, FormatResponse("error", "Huésped no encontrado", nil))
        return
    }

    log.Info().Str("id", id).Msg("Huésped actualizado exitosamente")
    c.JSON(http.StatusOK, FormatResponse("success", "Huésped actualizado correctamente", updatedGuest))
}