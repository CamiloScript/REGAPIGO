package APIController

import (
 
    "net/http"
  
    "github.com/CamiloScript/REGAPIGO/src/APIGuest/APIStruct"
    "github.com/CamiloScript/REGAPIGO/bd/MongoDB"
    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog/log"
    "go.mongodb.org/mongo-driver/bson"
 
)

// GetGuests maneja la solicitud GET para obtener todos los huéspedes
func GetGuests(c *gin.Context) {
    collection := MongoDB.Cliente.Database("APIREGDB").Collection("Guests")
    cursor, err := collection.Find(c, bson.M{})
    if err != nil {
        log.Error().Err(err).Msg("Error al obtener huéspedes desde MongoDB")
        c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al obtener los huéspedes", nil))
        return
    }
    defer cursor.Close(c)

    var guests []APIStruct.Guest
    if err := cursor.All(c, &guests); err != nil {
        log.Error().Err(err).Msg("Error al decodificar resultados de MongoDB")
        c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al procesar los huéspedes", nil))
        return
    }

    log.Info().Int("total_huéspedes", len(guests)).Msg("Huéspedes obtenidos exitosamente")
    c.JSON(http.StatusOK, FormatResponse("success", "Lista de huéspedes obtenida correctamente", guests))
}