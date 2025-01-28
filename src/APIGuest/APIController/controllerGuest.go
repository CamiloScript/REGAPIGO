package APIController


import (
	"github.com/CamiloScript/REGAPIGO/tree/main/src/APIGuest/APIStruct"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	

)

var Version = "1.0"

var Guests = []APIStruct.Guest{} // Simula una base de datos en memoria

func GetGuests(c *gin.Context) {
	log.Info().Msg("Fetching all guests")
	c.JSON(http.StatusOK, Guests)
}

func GetGuestByID(c *gin.Context) {
	id := c.Param("id")
	log.Info().Str("id", id).Msg("Fetching guest by ID")
	for _, guest := range Guests {
		if guest.ID == id {
			c.JSON(http.StatusOK, guest)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Guest not found"})
}

func CreateGuest(c *gin.Context) {
	var newGuest APIStruct.Guest
	if err := c.ShouldBindJSON(&newGuest); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	newGuest.ID = GenerateID() // Simula la generación de un ID único
	Guests = append(Guests, newGuest)
	log.Info().Str("id", newGuest.ID).Msg("Guest created successfully")
	c.JSON(http.StatusCreated, newGuest)
}

func UpdateGuest(c *gin.Context) {
	id := c.Param("id")
	var updatedData APIStruct.Guest
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	for i, guest := range Guests {
		if guest.ID == id {
			// Update fields selectively
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
			c.JSON(http.StatusOK, guest)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Guest not found"})
}

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
	c.JSON(http.StatusNotFound, gin.H{"error": "Guest not found"})
}

func GenerateID() string {
	// Simula la generación de un ID único
	return "some-unique-id"
}