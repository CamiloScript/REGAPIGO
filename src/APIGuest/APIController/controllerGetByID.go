package APIController

import (
    "net/http"
    "github.com/CamiloScript/REGAPIGO/src/APIGuest/APIStruct"
    "github.com/CamiloScript/REGAPIGO/bd/MongoDB"
    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog/log"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

// GetGuestByID maneja la solicitud GET para obtener un huésped por su ID
func GetGuestByID(c *gin.Context) {
    id := c.Param("id")
    log.Debug().Str("id_recibido", id).Msg("Iniciando búsqueda de huésped")

    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        log.Error().Err(err).Str("id", id).Msg("ID inválido")
        c.JSON(http.StatusBadRequest, FormatResponse("error", "ID inválido", nil))
        return
    }

    collection := MongoDB.Cliente.Database("APIREGDB").Collection("Guests")
    var guest APIStruct.Guest
    err = collection.FindOne(c, bson.M{"_id": objectID}).Decode(&guest)

    if err == mongo.ErrNoDocuments {
        log.Warn().Str("id", id).Msg("Huésped no encontrado en la base de datos")
        c.JSON(http.StatusNotFound, FormatResponse("error", "Huésped no encontrado", nil))
        return
    } else if err != nil {
        log.Error().Err(err).Str("id", id).Msg("Error al buscar huésped")
        c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al obtener el huésped", nil))
        return
    }

    log.Info().Str("id", id).Msg("Huésped encontrado exitosamente")
    c.JSON(http.StatusOK, FormatResponse("success", "Huésped encontrado", guest))
}
