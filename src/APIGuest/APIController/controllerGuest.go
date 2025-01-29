package APIController

import (
	"github.com/CamiloScript/REGAPIGO/src/APIGuest/APIStruct"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"github.com/CamiloScript/REGAPIGO/bd/MongoDB"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	client, err := ConectaMongoDB()
	if err != nil {
		log.Error().Err(err).Msg("Error al conectar a la base de datos")
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al conectar a la base de datos", nil))
		return
	}
	collection := client.Database("your_database").Collection("guests")
	cur, err := collection.Find(c, bson.D{})
	if err != nil {
		log.Error().Err(err).Msg("Error al obtener los huéspedes")
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "No se pudieron obtener los huéspedes", nil))
		return
	}
	defer cur.Close(c)

	var guests []APIStruct.Guest
	for cur.Next(c) {
		var guest APIStruct.Guest
		if err := cur.Decode(&guest); err != nil {
			log.Error().Err(err).Msg("Error al decodificar huésped")
			c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al procesar los huéspedes", nil))
			return
		}
		guests = append(guests, guest)
	}

	if err := cur.Err(); err != nil {
		log.Error().Err(err).Msg("Error al iterar sobre los huéspedes")
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al procesar los huéspedes", nil))
		return
	}

	c.JSON(http.StatusOK, FormatResponse("success", "Guests retrieved successfully", guests))
}

// Manejo de peticiones GET por ID para un huésped
func GetGuestByID(c *gin.Context) {
	id := c.Param("id") // Obtener ID de los parámetros de la URL
	objID, err := primitive.ObjectIDFromHex(id) // Convertir el ID de cadena a ObjectID de MongoDB
	if err != nil {
		log.Error().Err(err).Msg("Error al convertir ID")
		c.JSON(http.StatusBadRequest, FormatResponse("error", "ID inválido", nil))
		return
	}

	client, err := ConectaMongoDB()
	if err != nil {
		log.Error().Err(err).Msg("Error al conectar a la base de datos")
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al conectar a la base de datos", nil))
		return
	}
	collection := client.Database("your_database").Collection("guests")
	var guest APIStruct.Guest
	err = collection.FindOne(c, bson.M{"_id": objID}).Decode(&guest)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, FormatResponse("error", "Huésped no encontrado", nil))
		return
	} else if err != nil {
		log.Error().Err(err).Msg("Error al obtener el huésped")
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al obtener el huésped", nil))
		return
	}

	c.JSON(http.StatusOK, FormatResponse("success", "Guest retrieved successfully", guest))
}

// Manejo de peticiones POST para crear un nuevo huésped
func CreateGuest(c *gin.Context) {
	var guest APIStruct.Guest
	if err := c.ShouldBindJSON(&guest); err != nil {
		c.JSON(http.StatusBadRequest, FormatResponse("error", "Datos de huésped inválidos", nil))
		return
	}

	// Asignar fecha de creación
	guest.CreatedAt = time.Now()
	guest.UpdatedAt = time.Now()

	client, err := ConectaMongoDB()
	if err != nil {
		log.Error().Err(err).Msg("Error al conectar a la base de datos")
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al conectar a la base de datos", nil))
		return
	}
	collection := client.Database("your_database").Collection("guests")
	result, err := collection.InsertOne(c, guest)
	if err != nil {
		log.Error().Err(err).Msg("Error al crear el huésped")
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "No se pudo crear el huésped", nil))
		return
	}

	// Obtener ID generado por MongoDB
	guest.ID = result.InsertedID.(primitive.ObjectID).Hex()

	c.JSON(http.StatusCreated, FormatResponse("success", "Guest created successfully", guest))
}

// Manejo de peticiones PUT para actualizar un huésped
func UpdateGuest(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResponse("error", "ID inválido", nil))
		return
	}

	var guest APIStruct.Guest
	if err := c.ShouldBindJSON(&guest); err != nil {
		c.JSON(http.StatusBadRequest, FormatResponse("error", "Datos de huésped inválidos", nil))
		return
	}

	guest.UpdatedAt = time.Now()

	client, err := ConectaMongoDB()
	if err != nil {
		log.Error().Err(err).Msg("Error al conectar a la base de datos")
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al conectar a la base de datos", nil))
		return
	}
	collection := client.Database("your_database").Collection("guests")
	_, err = collection.UpdateOne(c, bson.M{"_id": objID}, bson.M{
		"$set": guest,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error al actualizar el huésped")
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "No se pudo actualizar el huésped", nil))
		return
	}

	c.JSON(http.StatusOK, FormatResponse("success", "Guest updated successfully", guest))
}

// Manejo de peticiones DELETE para eliminar un huésped
func DeleteGuest(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResponse("error", "ID inválido", nil))
		return
	}

	client, err := ConectaMongoDB()
	if err != nil {
		log.Error().Err(err).Msg("Error al conectar a la base de datos")
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al conectar a la base de datos", nil))
		return
	}
	collection := client.Database("your_database").Collection("guests")
	_, err = collection.DeleteOne(c, bson.M{"_id": objID})
	if err != nil {
		log.Error().Err(err).Msg("Error al eliminar el huésped")
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "No se pudo eliminar el huésped", nil))
		return
	}

	c.JSON(http.StatusOK, FormatResponse("success", "Guest deleted successfully", nil))
}

