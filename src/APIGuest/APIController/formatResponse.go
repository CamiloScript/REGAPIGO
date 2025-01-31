package APIController

import (
    "github.com/CamiloScript/REGAPIGO/src/APIGuest/APIStruct"
)

// La función FormatResponse se encarga de darle formato a la respuesta de la API, la cual es personalizable y escalable.
func FormatResponse(status, message string, data interface{}) APIStruct.ApiResponse {
    // Se puede modificar el formato de la respuesta de la API, según las necesidades del usuario, modificando la estructura de la respuesta.
    return APIStruct.ApiResponse{
        Status:  status,
        Message: message,
        Data:    data,
    }
}
