package APIController

import (
	"github.com/CamiloScript/REGAPIGO/src/APIGuest/APIStruct"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"github.com/CamiloScript/REGAPIGO/bd/MongoDB"
	"time"
)

// Estructura para respuestas estandarizadas
type ApiResponse struct {
	Time    string      `json:"time"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  string      `json:"status"`
}

// Version de la API
var Version = "1.0"

// Conexion a la base de datos MongoDB
var ConectaMongoDB = MongoDB.ConexionDB // variable para conectar realizar la conexion a la base de datos

// Lista de huéspedes (simulada por ahora)
var Guests = []APIStruct.Guest{} // slice de invitados

// Middleware para el manejo de logs
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		log.Info().Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Dur("duration", time.Since(start)).
			Msg("Request processed")
	}
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
	c.JSON(http.StatusOK, FormatResponse("success", "Guests retrieved successfully", Guests))
}

// Manejo de peticiones GET por ID de huésped
func GetGuestByID(c *gin.Context) {
	id := c.Param("id")
	log.Info().Str("id", id).Msg("Fetching guest by ID")
	for _, guest := range Guests {
		if guest.ID == id {
			c.JSON(http.StatusOK, FormatResponse("success", "Guest found", guest))
			return
		}
	}
	c.JSON(http.StatusNotFound, FormatResponse("error", "Guest not found", nil))
}

// Manejo de peticiones POST para crear un nuevo huésped
func CreateGuest(c *gin.Context) {
	var newGuest APIStruct.Guest
	if err := c.ShouldBindJSON(&newGuest); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, FormatResponse("error", "Invalid request body", nil))
		return
	}
	newGuest.ID = GenerateID() // Simula la generación de un ID único
	Guests = append(Guests, newGuest)
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

	for i, guest := range Guests {
		if guest.ID == id {
			// Actualiza campos selectivos
			if updatedData.FirstName != "" {
				guest.FirstName = updatedData.FirstName
			}
			if updatedData.LastName != "" {
				guest.LastName = updatedData.LastName
			}
			if updatedData.Email != "" {
				guest.Email = updatedData.Email
			}
			if updatedData.Address != (APIStruct.Address{}) {
				guest.Address = updatedData.Address
			}
			guest.Blacklisted = updatedData.Blacklisted
			guest.BlacklistReason = updatedData.BlacklistReason
			Guests[i] = guest
			log.Info().Str("id", id).Msg("Guest updated successfully")
			c.JSON(http.StatusOK, FormatResponse("success", "Guest updated successfully", guest))
			return
		}
	}
	c.JSON(http.StatusNotFound, FormatResponse("error", "Guest not found", nil))
}

// Manejo de peticiones DELETE para eliminar un huésped
func DeleteGuest(c *gin.Context) {
	id := c.Param("id")
	log.Info().Str("id", id).Msg("Attempting to delete guest")
	for i, guest := range Guests {
		if guest.ID == id {
			Guests = append(Guests[:i], Guests[i+1:]...)
			log.Info().Str("id", id).Msg("Guest deleted successfully")
			c.Status(http.StatusNoContent)
			return
		}
	}
	c.JSON(http.StatusNotFound, FormatResponse("error", "Guest not found", nil))
}

// Función para generar un ID único para cada huésped
func GenerateID() string {
	// Simula la generación de un ID único
	return "some-unique-id"
}
