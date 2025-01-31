package APIStruct

// ApiResponse representa la estructura de una respuesta de la API
type ApiResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

