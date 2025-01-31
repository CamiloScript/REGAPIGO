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

// UpdateGuest maneja la solicitud PUT para actualizar los datos de un huésped en la base de datos, a partir de su ID.
func UpdateGuest(c *gin.Context) {
    // Se obtiene el ID del huésped a actualizar, desde la solicitud.
    id := c.Param("id")
    // Se registra en el log que se está iniciando la actualización de un huésped.
    log.Debug().Str("id", id).Msg("Iniciando actualización de huésped")

    // Se convierte el ID a un tipo ObjectID de MongoDB, para realizar la búsqueda en la base de datos.
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        // Si el ID no es válido, se registra en el log y se responde con un error 400.
        log.Error().Err(err).Str("id", id).Msg("ID inválido para actualización")
        c.JSON(http.StatusBadRequest, FormatResponse("error", "ID inválido", nil))
        return
    }
    // var updatedGuest aloja los datos del huésped actualizados, recibidos en la solicitud.
    var updatedGuest APIStruct.Guest
    // Se comparan los datos recibidos en la solicitud con la estructura de un huésped.
    if err := c.ShouldBindJSON(&updatedGuest); err != nil {
        // Si los datos no cumplen con la estructura de un huésped, se registra en el log y se responde con un error 400.
        log.Error().Err(err).Msg("Datos inválidos en solicitud de actualización")
        c.JSON(http.StatusBadRequest, FormatResponse("error", "Datos inválidos", nil))
        return
    }
    // Se actualiza la fecha de actualización del huésped.
    updatedGuest.UpdatedAt = time.Now()

    // Se obtiene la colección de huéspedes en la base de datos.    
    collection := MongoDB.Cliente.Database("APIREGDB").Collection("Guests")
    // Se actualiza el huésped en la base de datos, a partir de su ID.
    result, err := collection.UpdateOne(c, bson.M{"_id": objectID}, bson.M{"$set": updatedGuest})
    if err != nil {
        // Si ocurre un error al actualizar el huésped, se registra en el log y se responde con un error 500.
        log.Error().Err(err).Str("id", id).Msg("Error al actualizar en MongoDB")
        c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al actualizar el huésped", nil))
        return
    }
    // Si no se encuentra un huésped con el ID proporcionado, se registra en el log y se responde con un error 404.
    if result.MatchedCount == 0 {
        log.Warn().Str("id", id).Msg("Huésped no encontrado para actualización")
        c.JSON(http.StatusNotFound, FormatResponse("error", "Huésped no encontrado", nil))
        return
    }
    // Se registra en el log que el huésped fue actualizado exitosamente.
    log.Info().Str("id", id).Msg("Huésped actualizado exitosamente")
    c.JSON(http.StatusOK, FormatResponse("success", "Huésped actualizado correctamente", updatedGuest))
}