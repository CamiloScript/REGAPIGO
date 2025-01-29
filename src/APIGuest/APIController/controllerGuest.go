package APIController

import (
	"github.com/CamiloScript/REGAPIGO/src/APIGuest/APIStruct"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"github.com/CamiloScript/REGAPIGO/bd/MongoDB"
	"go.mongodb.org/mongo-driver/bson"
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
	collection := ConectaMongoDB.Database("your_database").Collection("guests")
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
