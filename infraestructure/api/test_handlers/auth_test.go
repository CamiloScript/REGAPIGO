package test_handlers

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "bytes"
    "github.com/CamiloScript/REGAPIGO/infraestructure/api/handlers"
    "github.com/CamiloScript/REGAPIGO/infraestructure/persistence/servicio"
    "github.com/CamiloScript/REGAPIGO/shared/logger"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

// TestLoginCredencialesInvalidas - Prueba el caso de error cuando se envían credenciales inválidas.
// Este test verifica que el endpoint de login devuelva un error de formato inválido cuando no se proporcionan credenciales.
func TestLoginCredencialesInvalidas(t *testing.T) {
    // Inicializar el logger para las pruebas
    log := logger.NuevoRegistrador("TEST", "|")

    // Crear un mock del cliente de autenticación de Alfresco
    mockAuth := &servicio.MockAuthClient{Log: log}

    // Crear el manejador de autenticación utilizando el mock
    servicioAuth := handlers.NewAuthHandler(mockAuth, log)

    // Configurar el router de Gin para la prueba
    router := gin.Default()
    router.POST("/auth/login", servicioAuth.Login)

    // Crear una solicitud HTTP POST con un cuerpo vacío
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBufferString(``))
    req.Header.Set("Content-Type", "application/json")

    // Ejecutar la solicitud
    router.ServeHTTP(w, req)

    // Verificar que el código de estado HTTP sea 400 (Bad Request)
    assert.Equal(t, http.StatusBadRequest, w.Code)

    // Verificar que el cuerpo de la respuesta contenga el mensaje de error esperado
    assert.Contains(t, w.Body.String(), "Formato de JSON inválido")
}

// TestLoginErrorServicio - Prueba el caso de error cuando el servicio de autenticación falla.
// Este test simula un error en el servicio de autenticación y verifica que el endpoint de login devuelva un error 401 (Unauthorized).
func TestLoginErrorServicio(t *testing.T) {
    // Inicializar el logger para las pruebas
    log := logger.NuevoRegistrador("TEST", "|")

    // Crear un mock del cliente de autenticación de Alfresco
    mockAuth := &servicio.MockAuthClient{Log: log}

    // Crear el manejador de autenticación utilizando el mock
    servicioAuth := handlers.NewAuthHandler(mockAuth, log)

    // Forzar un error en el mock para simular un fallo en el servicio de autenticación
    mockAuth.ForzarError = true

    // Configurar el router de Gin para la prueba
    router := gin.Default()
    router.POST("/auth/login", servicioAuth.Login)

    // Crear una solicitud HTTP POST con credenciales inválidas
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBufferString(`{
        "userId": "admin",
        "password": "error"
    }`))
    req.Header.Set("Content-Type", "application/json")

    // Ejecutar la solicitud
    router.ServeHTTP(w, req)

    // Verificar que el código de estado HTTP sea 401 (Unauthorized)
    assert.Equal(t, http.StatusUnauthorized, w.Code)
}