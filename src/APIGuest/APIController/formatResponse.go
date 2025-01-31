package APIController

import (
    "github.com/CamiloScript/REGAPIGO/src/APIGuest/APIStruct"
)

// FormatResponse se mantiene igual
func FormatResponse(status, message string, data interface{}) APIStruct.ApiResponse {
    return APIStruct.ApiResponse{
        Status:  status,
        Message: message,
        Data:    data,
    }
}
