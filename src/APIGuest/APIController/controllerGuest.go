package APIController

import (
    "fmt"
    "net/http"
    "time"
    
    "github.com/CamiloScript/REGAPIGO/src/APIGuest/APIStruct"
    "github.com/CamiloScript/REGAPIGO/bd/MongoDB"
    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog/log"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
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

// CreateGuest maneja la solicitud POST para crear un nuevo huésped
func CreateGuest(c *gin.Context) {
    var newGuest APIStruct.Guest
    if err := c.ShouldBindJSON(&newGuest); err != nil {
        log.Error().Err(err).Msg("Datos inválidos en solicitud de creación")
        c.JSON(http.StatusBadRequest, FormatResponse("error", "Datos inválidos", nil))
        return
    }

    newGuest.CreatedAt = time.Now()
    newGuest.UpdatedAt = time.Now()

    collection := MongoDB.Cliente.Database("APIREGDB").Collection("Guests")
    result, err := collection.InsertOne(c, newGuest)
    if err != nil {
        log.Error().Err(err).Interface("huésped", newGuest).Msg("Error al insertar en MongoDB")
        c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al crear el huésped", nil))
        return
    }

    log.Info().
        Str("id_insertado", fmt.Sprintf("%v", result.InsertedID)).
        Msg("Nuevo huésped creado")
    
    c.JSON(http.StatusCreated, FormatResponse("success", "Huésped creado correctamente", newGuest))
}

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

// FormatResponse se mantiene igual
func FormatResponse(status, message string, data interface{}) APIStruct.ApiResponse {
    return APIStruct.ApiResponse{
        Status:  status,
        Message: message,
        Data:    data,
    }
}
