package APIController_test

import (
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/CamiloScript/REGAPIGO/src/APIGuest/APIController"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
)

func TestGetGuests(t *testing.T) {
	router := gin.Default()
	router.GET("/guests", APIController.GetGuests)

	req, _ := http.NewRequest("GET", "/guests", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Guests retrieved successfully")
}

func TestCreateGuest(t *testing.T) {
	router := gin.Default()
	router.POST("/guests", APIController.CreateGuest)

	req, _ := http.NewRequest("POST", "/guests", nil) // Añadir cuerpo JSON de ejemplo aquí
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
}
