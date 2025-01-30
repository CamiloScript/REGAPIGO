package APIController_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/CamiloScript/REGAPIGO/src/APIGuest/APIController"
	"github.com/CamiloScript/REGAPIGO/src/APIGuest/APIStruct"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

// Configuración inicial para zerolog
func init() {
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
    log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func TestGetGuests(t *testing.T) {
    // Configurar el router de Gin
    router := gin.Default()
    router.GET("/guests", APIController.GetGuests)

    // Crear una solicitud HTTP GET
    req, err := http.NewRequest("GET", "/guests", nil)
    if err != nil {
        log.Fatal().Err(err).Msg("Error al crear la solicitud HTTP para GetGuests")
    }

    // Registrar la prueba
    log.Info().Str("endpoint", "/guests").Msg("Iniciando prueba de GetGuests")

    // Ejecutar la solicitud
    resp := httptest.NewRecorder()
    router.ServeHTTP(resp, req)

    // Verificar el código de estado
    assert.Equal(t, http.StatusOK, resp.Code, "El código de estado debería ser 200")

    // Verificar el cuerpo de la respuesta
    var response APIStruct.ApiResponse
    err = json.Unmarshal(resp.Body.Bytes(), &response)
    if err != nil {
        log.Error().Err(err).Msg("Error al decodificar la respuesta JSON")
        t.FailNow()
    }

    assert.Equal(t, "success", response.Status, "El estado de la respuesta debería ser 'success'")
    assert.Contains(t, response.Message, "Lista de huéspedes obtenida correctamente", "El mensaje de respuesta es incorrecto")

    log.Info().Str("status", response.Status).Msg("Prueba de GetGuests completada exitosamente")
}

func TestCreateGuest(t *testing.T) {
    // Configurar el router de Gin
    router := gin.Default()
    router.POST("/guests", APIController.CreateGuest)

    // Crear un huésped de prueba
    newGuest := APIStruct.Guest{
        FirstName: "Jane",
        LastName:  "Doe",
        Email:     "jane.doe@example.com",
    }
    guestJSON, err := json.Marshal(newGuest)
    if err != nil {
        log.Fatal().Err(err).Msg("Error al convertir el huésped a JSON")
    }

    // Crear una solicitud HTTP POST con el cuerpo JSON
    req, err := http.NewRequest("POST", "/guests", bytes.NewBuffer(guestJSON))
    if err != nil {
        log.Fatal().Err(err).Msg("Error al crear la solicitud HTTP para CreateGuest")
    }
    req.Header.Set("Content-Type", "application/json")

    // Registrar la prueba
    log.Info().Str("endpoint", "/guests").Msg("Iniciando prueba de CreateGuest")

    // Ejecutar la solicitud
    resp := httptest.NewRecorder()
    router.ServeHTTP(resp, req)

    // Verificar el código de estado
    assert.Equal(t, http.StatusCreated, resp.Code, "El código de estado debería ser 201")

    // Verificar el cuerpo de la respuesta
    var response APIStruct.ApiResponse
    err = json.Unmarshal(resp.Body.Bytes(), &response)
    if err != nil {
        log.Error().Err(err).Msg("Error al decodificar la respuesta JSON")
        t.FailNow()
    }

    assert.Equal(t, "success", response.Status, "El estado de la respuesta debería ser 'success'")
    assert.Contains(t, response.Message, "Huésped creado correctamente", "El mensaje de respuesta es incorrecto")

    log.Info().Str("status", response.Status).Msg("Prueba de CreateGuest completada exitosamente")
}