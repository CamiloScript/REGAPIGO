package auth



import (
    "github.com/CamiloScript/REGAPIGO/infraestructure/persistence/servicio" // Cliente de autenticación para Alfresco.
    "github.com/CamiloScript/REGAPIGO/shared/logger" // Módulo de logging centralizado.
    "fmt" // Formateo de strings.
)

// AuthService define la interfaz para la autenticación de usuarios.
// Esta interfaz facilita la abstracción y permite cambiar la implementación sin afectar el resto del código.
type AuthService interface {
    Authenticate(usuario, password string) (string, error) // Retorna un token de autenticación o un error.
}

// AuthServiceImpl es la implementación concreta del servicio de autenticación.
// Se conecta con el cliente de autenticación de Alfresco y gestiona los registros de eventos.
type AuthServiceImpl struct {
    client servicio.AuthClient // Cliente que maneja la autenticación con el servicio externo Alfresco.
    log    *logger.Registrador // Logger para registrar eventos de autenticación.
}

// NewAuthService es un constructor que inicializa AuthServiceImpl con sus dependencias.
// Utiliza inyección de dependencias para permitir flexibilidad y pruebas más sencillas.
func NewAuthService(client servicio.AuthClient, log *logger.Registrador) *AuthServiceImpl {
    return &AuthServiceImpl{client: client, log: log}
}

// Authenticate valida las credenciales del usuario y delega la autenticación al cliente de Alfresco.
// - `usuario`: Nombre de usuario.
// - `password`: Contraseña del usuario.
// Retorna:
//   - `string`: Token de autenticación en caso de éxito.
//   - `error`: Error si la autenticación falla.
func (s *AuthServiceImpl) Authenticate(usuario, password string) (string, error) {
    return s.client.Login(usuario, password)
}

// MockAuthClient es tu cliente de autenticación simulado
type MockAuthClient struct {
    Log          logger.Registrador
    ForzarError  bool
}

// Implementa el método Authenticate de la interfaz auth.AuthService
func (m *MockAuthClient) Authenticate(userID, password string) (string, error) {
    if m.ForzarError {
        return "", fmt.Errorf("error de autenticación")
    }
    // Aquí va la lógica para simular la autenticación
    return "token_simulado", nil
}