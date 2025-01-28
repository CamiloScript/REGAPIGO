package APIStruct

import "time"

type Address struct {
	Street     string `json:"street,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	Country    string `json:"country,omitempty"`
	PostalCode string `json:"postalCode,omitempty"`
}

// Guest struct para el huésped
type Guest struct {
	ID             string            `json:"id"`                                 // ID único (puede ser generado por una base de datos)
	FirstName      string            `json:"firstName" binding:"required"`       // Requerido
	LastName       string            `json:"lastName" binding:"required"`        // Requerido
	Email          string            `json:"email" binding:"required,email"`     // Requerido y validación de formato
	Phone          string            `json:"phone,omitempty"`                    // Opcional
	Nationality    string            `json:"nationality,omitempty"`              // Opcional
	DocumentType   string            `json:"documentType,omitempty"`             // Opcional, enum en validación
	DocumentNumber string            `json:"documentNumber,omitempty"`           // Opcional
	Address        Address           `json:"address,omitempty"`                  // Estructura anidada
	Blacklisted    bool              `json:"blacklisted"`                        // Booleano con valor predeterminado
	BlacklistReason string           `json:"blacklistReason,omitempty"`          // Opcional
	CreatedAt      time.Time         `json:"createdAt"`                          // Marca de tiempo
	UpdatedAt      time.Time         `json:"updatedAt"`                          // Marca de tiempo
}