package APIController

import (
 
    "net/http"
  
    "github.com/CamiloScript/REGAPIGO/src/APIGuest/APIStruct"
    "github.com/CamiloScript/REGAPIGO/bd/MongoDB"
    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog/log"
    "go.mongodb.org/mongo-driver/bson"
 
)

// GetGuests maneja la solicitud GET para obtener todos los huéspedes presentes en la coleccion Guests de la base de datos.
func GetGuests(c *gin.Context) {
    // Se obtienen todos los huéspedes presentes en la base de datos, y se registra con zerolog para su seguimiento.
    collection := MongoDB.Cliente.Database("APIREGDB").Collection("Guests")
    cursor, err := collection.Find(c, bson.M{})
    // Se registran los eventos con zerolog si ocurre un error al obtener los huéspedes, y se retorna un mensaje de error.
    if err != nil {
        log.Error().Err(err).Msg("Error al obtener huéspedes desde MongoDB")
        c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al obtener los huéspedes", nil))
        return
    }
    // Se cierra el cursor correspondiente a la busqueda en la base de datos.
    defer cursor.Close(c)

    // Se decodifican los resultados obtenidos de la base de datos, y se registra con zerolog para su seguimiento.
    var guests []APIStruct.Guest
    // Se registran los eventos con zerolog si ocurre un error al decodificar los resultados, y se retorna un mensaje de error.
    if err := cursor.All(c, &guests); err != nil {
        log.Error().Err(err).Msg("Error al decodificar resultados de MongoDB")
        c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al procesar los huéspedes", nil))
        return
    }
    // Se registran los eventos con zerolog si los huéspedes son obtenidos exitosamente, y se retorna un mensaje de éxito.
    // Adicional al mensaje de exito se adjunta la lista de huéspedes obtenidos.
    log.Info().Int("total_huéspedes", len(guests)).Msg("Huéspedes obtenidos exitosamente")
    c.JSON(http.StatusOK, FormatResponse("success", "Lista de huéspedes obtenida correctamente", guests))
}