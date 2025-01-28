package APIController

import (
	"github.com/CamiloScript/REGAPIGO/src/APIGuest/APIStruct"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"github.com/CamiloScript/REGAPIGO/bd/MongoDB"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var ConectaMongoDB = MongoDB.ConexionDB // Variable para conectar a la base de datos

// Estructura para respuestas estandarizadas
type ApiResponse struct {
	Time    string      `json:"time"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  string      `json:"status"`
}

// Función para formatear respuestas
func FormatResponse(status string, message string, data interface{}) ApiResponse {
	return ApiResponse{
		Time:    time.Now().Format(time.RFC3339),
		Status:  status,
		Message: message,
		Data:    data,
	}
}

// Manejo de peticiones GET para todos los huéspedes
func GetGuests(c *gin.Context) {
	log.Info().Msg("Fetching all guests")

	// Consultar todos los huéspedes en la base de datos
	collection := ConectaMongoDB.Database("your_db").Collection("guests")
	cursor, err := collection.Find(c, bson.M{})
	if err != nil {
		log.Error().Err(err).Msg("Error fetching guests")
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error fetching guests", nil))
		return
	}
	defer cursor.Close(c)

	var guests []APIStruct.Guest
	if err = cursor.All(c, &guests); err != nil {
		log.Error().Err(err).Msg("Error parsing guests data")
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error parsing guests data", nil))
		return
	}

	c.JSON(http.StatusOK, FormatResponse("success", "Guests retrieved successfully", guests))
}

// Manejo de peticiones GET por ID de huésped
func GetGuestByID(c *gin.Context) {
	id := c.Param("id")
	log.Info().Str("id", id).Msg("Fetching guest by ID")

	// Buscar el huésped por ID en la base de datos
	collection := ConectaMongoDB.Database("your_db").Collection("guests")
	var guest APIStruct.Guest
	err := collection.FindOne(c, bson.M{"id": id}).Decode(&guest)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, FormatResponse("error", "Guest not found", nil))
		} else {
			log.Error().Err(err).Msg("Error fetching guest")
			c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error fetching guest", nil))
		}
		return
	}

	c.JSON(http.StatusOK, FormatResponse("success", "Guest found", guest))
}

// Manejo de peticiones POST para crear un nuevo huésped
func CreateGuest(c *gin.Context) {
	var newGuest APIStruct.Guest
	if err := c.ShouldBindJSON(&newGuest); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, FormatResponse("error", "Invalid request body", nil))
		return
	}

	// Conectar a la base de datos para insertar un nuevo huésped
	collection := ConectaMongoDB.Database("your_db").Collection("guests")
	newGuest.ID = GenerateID() // Generar un ID único para el nuevo huésped

	_, err := collection.InsertOne(c, newGuest)
	if err != nil {
		log.Error().Err(err).Msg("Error creating guest")
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error creating guest", nil))
		return
	}

	log.Info().Str("id", newGuest.ID).Msg("Guest created successfully")
	c.JSON(http.StatusCreated, FormatResponse("success", "Guest created successfully", newGuest))
}

// Manejo de peticiones PUT para actualizar un huésped
func UpdateGuest(c *gin.Context) {
	id := c.Param("id")
	var updatedData APIStruct.Guest
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, FormatResponse("error", "Invalid request body", nil))
		return
	}

	// Conectar a la base de datos y actualizar el huésped por ID
	collection := ConectaMongoDB.Database("your_db").Collection("guests")
	_, err := collection.UpdateOne(c, bson.M{"id": id}, bson.M{
		"$set": updatedData,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error updating guest")
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error updating guest", nil))
		return
	}

	log.Info().Str("id", id).Msg("Guest updated successfully")
	c.JSON(http.StatusOK, FormatResponse("success", "Guest updated successfully", updatedData))
}

// Manejo de peticiones DELETE para eliminar un huésped
func DeleteGuest(c *gin.Context) {
	id := c.Param("id")
	log.Info().Str("id", id).Msg("Attempting to delete guest")

	// Buscar y eliminar el huésped de la base de datos
	collection := ConectaMongoDB.Database("your_db").Collection("guests")
	_, err := collection.DeleteOne(c, bson.M{"id": id})
	if err != nil {
		log.Error().Err(err).Msg("Error deleting guest")
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error deleting guest", nil))
		return
	}

	log.Info().Str("id", id).Msg("Guest deleted successfully")
	c.Status(http.StatusNoContent)
}

// Función para generar un ID único para cada huésped
func GenerateID() string {
	// Simula la generación de un ID único
	return "some-unique-id"
}
