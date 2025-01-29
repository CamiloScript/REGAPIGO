package APIController

import (
	"net/http"

	"github.com/CamiloScript/REGAPIGO/src/models"
	"github.com/CamiloScript/REGAPIGO/bd/MongoDB"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetGuests maneja la solicitud GET para obtener todos los huéspedes
func GetGuests(c *gin.Context) {
	// Consultamos la base de datos
	collection := MongoDB.Client.Database("guestsDB").Collection("guests")
	cursor, err := collection.Find(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al obtener los huéspedes", nil))
		return
	}
	defer cursor.Close(c)

	var guests []models.Guest
	if err := cursor.All(c, &guests); err != nil {
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al procesar los huéspedes", nil))
		return
	}

	c.JSON(http.StatusOK, FormatResponse("success", "Lista de huéspedes obtenida correctamente", guests))
}

// GetGuestByID maneja la solicitud GET para obtener un huésped por su ID
func GetGuestByID(c *gin.Context) {
	id := c.Param("id")

	// Consultamos la base de datos
	collection := MongoDB.Client.Database("guestsDB").Collection("guests")
	var guest models.Guest
	err := collection.FindOne(c, bson.M{"_id": id}).Decode(&guest)

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, FormatResponse("error", "Huésped no encontrado", nil))
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al obtener el huésped", nil))
		return
	}

	c.JSON(http.StatusOK, FormatResponse("success", "Huésped encontrado", guest))
}

// CreateGuest maneja la solicitud POST para crear un nuevo huésped
func CreateGuest(c *gin.Context) {
	var newGuest models.Guest
	if err := c.ShouldBindJSON(&newGuest); err != nil {
		c.JSON(http.StatusBadRequest, FormatResponse("error", "Datos inválidos", nil))
		return
	}

	// Insertamos el nuevo huésped en la base de datos
	collection := MongoDB.Client.Database("guestsDB").Collection("guests")
	_, err := collection.InsertOne(c, newGuest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al crear el huésped", nil))
		return
	}

	c.JSON(http.StatusCreated, FormatResponse("success", "Huésped creado correctamente", newGuest))
}

// UpdateGuest maneja la solicitud PUT para actualizar los datos de un huésped
func UpdateGuest(c *gin.Context) {
	id := c.Param("id")
	var updatedGuest models.Guest
	if err := c.ShouldBindJSON(&updatedGuest); err != nil {
		c.JSON(http.StatusBadRequest, FormatResponse("error", "Datos inválidos", nil))
		return
	}

	// Actualizamos el huésped en la base de datos
	collection := MongoDB.Client.Database("guestsDB").Collection("guests")
	_, err := collection.UpdateOne(c, bson.M{"_id": id}, bson.M{"$set": updatedGuest})
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, FormatResponse("error", "Huésped no encontrado", nil))
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al actualizar el huésped", nil))
		return
	}

	c.JSON(http.StatusOK, FormatResponse("success", "Huésped actualizado correctamente", updatedGuest))
}

// DeleteGuest maneja la solicitud DELETE para eliminar un huésped
func DeleteGuest(c *gin.Context) {
	id := c.Param("id")

	// Eliminamos el huésped de la base de datos
	collection := MongoDB.Client.Database("guestsDB").Collection("guests")
	_, err := collection.DeleteOne(c, bson.M{"_id": id})
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, FormatResponse("error", "Huésped no encontrado", nil))
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al eliminar el huésped", nil))
		return
	}

	c.JSON(http.StatusOK, FormatResponse("success", "Huésped eliminado correctamente", nil))
}

// FormatResponse formatea las respuestas de la API
func FormatResponse(status, message string, data interface{}) gin.H {
	return gin.H{
		"status":  status,
		"message": message,
		"data":    data,
	}
}
