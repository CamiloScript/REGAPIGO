package handlers

import (
    "net/http"
    "github.com/CamiloScript/REGAPIGO/domain/auth"
    "github.com/CamiloScript/REGAPIGO/shared/logger"
    "github.com/gin-gonic/gin"
)

// AuthHandler maneja las operaciones de autenticación sin sesiones.
// Este manejador es responsable de procesar las solicitudes de autenticación
// y devolver un ticket de autenticación directamente al cliente.
type AuthHandler struct {
    servicio auth.AuthService // Servicio de autenticación inyectado
    log      *logger.Registrador // Logger para registro de eventos
}

// NewAuthHandler inicializa el manejador con dependencias necesarias.
// Recibe un servicio de autenticación y un logger para registrar eventos.
// Retorna una instancia de AuthHandler lista para ser utilizada.
func NewAuthHandler(servicio auth.AuthService, log *logger.Registrador) *AuthHandler {
    return &AuthHandler{servicio: servicio, log: log}
}

// Login procesa la autenticación y devuelve el ticket directamente.
// Este método maneja la solicitud POST para autenticar a un usuario.
// Si las credenciales son válidas, devuelve un ticket de autenticación.
// Si las credenciales son inválidas o el formato de la solicitud es incorrecto,
// devuelve un error correspondiente.
func (h *AuthHandler) Login(c *gin.Context) {
    // Estructura para capturar credenciales del cuerpo de la solicitud
    var credenciales struct {
        UserId   string `json:"userId"`   // Nombre de usuario
        Password string `json:"password"` // Contraseña
    }

    // Vincular JSON a la estructura
    // Si el JSON no coincide con la estructura o es inválido, se devuelve un error 400 (Bad Request).
    if err := c.ShouldBindJSON(&credenciales); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de JSON inválido"})
        return
    }

    if credenciales.UserId == "" || credenciales.Password == "" {
        h.log.Error("Credenciales faltantes", nil)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Credenciales faltantes"})
        return
    }

    // Autenticar con el servicio de autenticación
    // Se llama al servicio de autenticación para validar las credenciales.
    // Si el servicio devuelve un error, se registra el error y se devuelve un error 401 (Unauthorized).
    ticket, err := h.servicio.Authenticate(credenciales.UserId, credenciales.Password)
    if err != nil {
        h.log.Error("Credenciales inválidas", map[string]interface{}{"usuario": credenciales.UserId})
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
        return
    }

    // Respuesta exitosa con ticket (sin almacenar en sesión)
    // Si la autenticación es exitosa, se devuelve un código 200 (OK) junto con el ticket.
    c.JSON(http.StatusOK, gin.H{"ticket": ticket})
}