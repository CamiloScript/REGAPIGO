package APIStruct

import (

	 "time"
	 "go.mongodb.org/mongo-driver/bson/primitive"
)

// Address representa la dirección de un huésped
type Address struct {
	Street     string `json:"street,omitempty"`     // Calle
	City       string `json:"city,omitempty"`       // Ciudad
	State      string `json:"state,omitempty"`      // Estado
	Country    string `json:"country,omitempty"`    // País
	PostalCode string `json:"postalCode,omitempty"` // Código postal
}

// Guest representa la información de un huésped
type Guest struct {
	ID 			   primitive.ObjectID	 `json:"id" bson:"_id,omitempty"`			  // ID del huesped generado por la base de datos al guardarlo
	FirstName      string    			 `json:"firstName" binding:"required"`        // Primer nombre (requerido)
	LastName       string    			 `json:"lastName" binding:"required"`         // Apellido (requerido)
	Email          string    			 `json:"email" binding:"required,email"`      // Correo electrónico (requerido y validación de formato)
	Phone          string    			 `json:"phone,omitempty"`                     // Teléfono (opcional)
	Nationality    string    			 `json:"nationality,omitempty"`               // Nacionalidad (opcional)
	DocumentType   string   	 	     `json:"documentType,omitempty"`              // Tipo de documento (opcional)
	DocumentNumber string    			 `json:"documentNumber,omitempty"`            // Número de documento (opcional)
	Address        Address  			 `json:"address,omitempty"`                   // Dirección (estructura anidada, opcional)
	Blacklisted    bool      			 `json:"blacklisted"`                         // Si está en lista negra (requerido)
	BlacklistReason string  			 `json:"blacklistReason,omitempty"`           // Motivo de la lista negra (opcional)
	CreatedAt      time.Time 			 `json:"createdAt"`                           // Fecha de creación
	UpdatedAt      time.Time 			 `json:"updatedAt"`                           // Fecha de última actualización
}

type ApiResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

