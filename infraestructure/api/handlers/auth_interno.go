package handlers

import (
    "fmt"
    "github.com/CamiloScript/REGAPIGO/domain/auth"
    "github.com/CamiloScript/REGAPIGO/shared/logger"
    "github.com/CamiloScript/REGAPIGO/shared/config"
)

// InternalAuth es una estructura que encapsula la lógica de autenticación interna.
type InternalAuth struct {
    authServicio auth.AuthService // Servicio de autenticación
    log          *logger.Registrador // Logger para registro de eventos
    cfg          *config.Config // Configuración de la aplicación
}

// NewInternalAuth crea una nueva instancia de InternalAuth.
func NewInternalAuth(authServicio auth.AuthService, log *logger.Registrador, cfg *config.Config) *InternalAuth {
    return &InternalAuth{
        authServicio: authServicio,
        log:          log,
        cfg:          cfg,
    }
}

// AutenticarInternamente genera un ticket de autenticación con Alfresco.
// Retorna el ticket y un error en caso de fallo.
func (ia *InternalAuth) AutenticarInternamente() (string, error) {
    // Obtener credenciales desde la configuración
    user := ia.cfg.AuthUser
    password := ia.cfg.AuthPassword

    // Llamar al servicio de autenticación
    ticket, err := ia.authServicio.Authenticate(user, password)
    if err != nil {
        ia.log.Error("Error de autenticación interna", map[string]interface{}{
            "error": err.Error(),
            "user":  user,
        })
        return "", fmt.Errorf("error al autenticar internamente: %v", err)
    }

    ia.log.Info("Autenticación interna exitosa", map[string]interface{}{
        "user": user,
    })
    return ticket, nil
}