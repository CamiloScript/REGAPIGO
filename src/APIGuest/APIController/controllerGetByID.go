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

// La función GetGuestByID maneja la solicitud GETBYID para obtener un huésped por su ID
func GetGuestByID(c *gin.Context) {

    // Se obtiene el ID del huésped a buscar desde la solicitud, y se registra con zerolog para su seguimiento.
    id := c.Param("id")
    // Se registra con zerolog el inicio de la busqueda del huésped.
    log.Debug().Str("id_recibido", id).Msg("Iniciando búsqueda de huésped")

    // Se transforma el ID del huésped a un ObjectID de MongoDB, para poder buscarlo en la base de datos y se registra con zerolog para su seguimiento.
    objectID, err := primitive.ObjectIDFromHex(id)
    // Se registra con zerolog si el ID del huésped es inválido, y se retorna un mensaje de error.
    if err != nil {
        log.Error().Err(err).Str("id", id).Msg("ID inválido")
        c.JSON(http.StatusBadRequest, FormatResponse("error", "ID inválido", nil))
        return
    }
    // Se realiza la busqueda del id del huésped en la base de datos, y se registra con zerolog para su seguimiento.
    collection := MongoDB.Cliente.Database("APIREGDB").Collection("Guests")
    // Se genera la variable guest para almacenar la estructura del huésped y completarla con los datos obtenidos de la base de datos.
    var guest APIStruct.Guest
    // Se realiza la busqueda del huésped en la base de datos, y se registra con zerolog para su seguimiento.
    err = collection.FindOne(c, bson.M{"_id": objectID}).Decode(&guest)
    // Se registra con zerolog si el huésped no es encontrado en la base de datos, y se retorna un mensaje de error.
    if err == mongo.ErrNoDocuments {
        log.Warn().Str("id", id).Msg("Huésped no encontrado en la base de datos")
        c.JSON(http.StatusNotFound, FormatResponse("error", "Huésped no encontrado", nil))
        return
    // Se registra con zerolog si ocurre un error al buscar el huésped, y se retorna un mensaje de error.
    } else if err != nil {
        log.Error().Err(err).Str("id", id).Msg("Error al buscar huésped")
        c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al obtener el huésped", nil))
        return
    }
    // Se registra con zerolog si el huésped es encontrado exitosamente, y se retorna un mensaje de éxito.
    log.Info().Str("id", id).Msg("Huésped encontrado exitosamente")
    c.JSON(http.StatusOK, FormatResponse("success", "Huésped encontrado", guest))
}
