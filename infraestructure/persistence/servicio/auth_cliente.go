package servicio

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "github.com/CamiloScript/REGAPIGO/shared/config"
    "github.com/CamiloScript/REGAPIGO/shared/logger"
)

// AuthClientImpl implementa la autenticación con Alfresco.
type AuthClientImpl struct {
    baseURL  string              // URL base del servidor Alfresco
    apiKey   string              // API Key para autenticación
    client   *http.Client        // Cliente HTTP
    log      *logger.Registrador // Logger
}

// AuthClient define la interfaz para la autenticación.
type AuthClient interface {
    Login(usuario, password string) (string, error)
}

// LoginResponse representa la estructura de respuesta de Alfresco.
type LoginResponse struct {
    ID     string   `json:"id"`     // Ticket de autenticación
    Roles  []string `json:"roles"`  // Roles del usuario (no usado actualmente)
    UserId string   `json:"userId"` // ID del usuario (no usado actualmente)
}

// NewAuthClient crea un cliente de autenticación configurado.
// Parámetros:
//   - cfg: Configuración de la aplicación.
//   - log: Logger para registrar eventos y errores.
// Retorna un cliente de autenticación configurado.
func NewAuthClient(cfg *config.Config, log *logger.Registrador) AuthClient {
    return &AuthClientImpl{
        baseURL: cfg.AlfrescoBaseURL,
        apiKey:  cfg.AlfrescoAPIKey,
        client: &http.Client{
            Timeout: 20 * time.Second, // Timeout para evitar bloqueos
        },
        log: log,
    }
}

// Login realiza la autenticación con Alfresco incluyendo la API Key.
// Parámetros:
//   - usuario: Nombre de usuario para autenticación.
//   - password: Contraseña del usuario.
// Retorna el ticket de autenticación o un error en caso de fallo.
func (c *AuthClientImpl) Login(usuario, password string) (string, error) {
    endpoint := "/session/log-in"
    url := c.baseURL + endpoint

    // 1. Registrar inicio de autenticación
    c.log.Info("Iniciando autenticación", map[string]interface{}{
        "usuario": usuario,
        "url":     url,
    })

    // 2. Crear cuerpo de la solicitud
    body := map[string]string{
        "userId":   usuario,
        "password": password,
    }
    jsonBody, err := json.Marshal(body)
    if err != nil {
        c.log.Error("Error al serializar cuerpo", map[string]interface{}{"error": err.Error()})
        return "", fmt.Errorf("error en formato de solicitud: %v", err)
    }

    // 3. Crear solicitud HTTP
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
    if err != nil {
        c.log.Error("Error al crear solicitud", map[string]interface{}{"error": err.Error()})
        return "", fmt.Errorf("error de conexión: %v", err)
    }

    // 4. Configurar headers requeridos
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("ADFTannerServices", c.apiKey) // API Key desde configuración

    // 5. Enviar solicitud
    resp, err := c.client.Do(req)
    if err != nil {
        c.log.Error("Error de red", map[string]interface{}{"error": err.Error()})
        return "", fmt.Errorf("error de comunicación: %v", err)
    }
    defer resp.Body.Close()

    // 6. Manejar errores HTTP
    if resp.StatusCode != http.StatusOK {
        c.log.Error("Error de autenticación", map[string]interface{}{
            "status_code": resp.StatusCode,
            "status":      resp.Status,
        })
        return "", fmt.Errorf("error de autenticación: %s", resp.Status)
    }

    // 7. Decodificar respuesta
    var loginResp LoginResponse
    if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
        c.log.Error("Error al decodificar respuesta", map[string]interface{}{"error": err.Error()})
        return "", fmt.Errorf("error en formato de respuesta: %v", err)
    }

    // 8. Registrar éxito
    c.log.Info("Autenticación exitosa", map[string]interface{}{"ticket": loginResp.ID})
    return loginResp.ID, nil
}

