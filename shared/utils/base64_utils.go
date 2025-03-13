package utils

import (
    "bytes"
    "encoding/base64"
    "errors"
    "strings"
)

// MimeTypes soportados y sus firmas base64
var mimeTypeSignatures = map[string]string{
    "image/jpeg": "/9j/",
    "image/png":  "iVBORw0KGgo",
    "application/pdf": "JVBERi0",
}

// DetectContentType detecta el tipo MIME basado en la firma del archivo base64
func DetectContentType(base64Data string) (string, error) {
    // Eliminar el prefijo "data:image/png;base64," si existe
    cleanData := base64Data
    if strings.Contains(base64Data, ";base64,") {
        parts := strings.Split(base64Data, ";base64,")
        if len(parts) == 2 {
            cleanData = parts[1]
        }
    }

    for mimeType, signature := range mimeTypeSignatures {
        if strings.HasPrefix(cleanData, signature) {
            return mimeType, nil
        }
    }
    return "", errors.New("formato de archivo no reconocido o no soportado")
}

// GetFileExtension obtiene la extensión del archivo basada en el tipo MIME
func GetFileExtension(mimeType string) string {
    switch mimeType {
    case "image/jpeg":
        return "jpg"
    case "image/png":
        return "png"
    case "application/pdf":
        return "pdf"
    default:
        return ""
    }
}

// IsValidBase64 verifica si el string proporcionado es un base64 válido
func IsValidBase64(base64String string) bool {
    // Eliminar el prefijo si existe
    cleanData := base64String
    if strings.Contains(base64String, ";base64,") {
        parts := strings.Split(base64String, ";base64,")
        if len(parts) == 2 {
            cleanData = parts[1]
        }
    }

    // Verificar si se puede decodificar
    _, err := base64.StdEncoding.DecodeString(cleanData)
    return err == nil
}

// DecodeBase64 decodifica un string base64 a bytes
func DecodeBase64(base64String string) ([]byte, error) {
    // Eliminar el prefijo si existe
    cleanData := base64String
    if strings.Contains(base64String, ";base64,") {
        parts := strings.Split(base64String, ";base64,")
        if len(parts) == 2 {
            cleanData = parts[1]
        }
    }

    return base64.StdEncoding.DecodeString(cleanData)
}

// EncodeToBase64 codifica bytes a string base64
func EncodeToBase64(data []byte) string {
    return base64.StdEncoding.EncodeToString(data)
}

// ConvertBase64ToMultipartFile convierte datos base64 a un buffer que puede usarse en lugar de un archivo
func ConvertBase64ToBuffer(base64Data string) (*bytes.Buffer, error) {
    fileBytes, err := DecodeBase64(base64Data)
    if err != nil {
        return nil, err
    }
    
    return bytes.NewBuffer(fileBytes), nil
}

// DetectMimeTypeFromContent determina el tipo MIME basado en el contenido del archivo
func DetectMimeTypeFromContent(content []byte) string {
    if len(content) > 4 {
        // Detectar PDF
        if bytes.HasPrefix(content, []byte("%PDF")) {
            return "application/pdf"
        }
        
        // Detectar JPEG
        if bytes.HasPrefix(content, []byte{0xFF, 0xD8, 0xFF}) {
            return "image/jpeg"
        }
        
        // Detectar PNG
        if bytes.HasPrefix(content, []byte{0x89, 0x50, 0x4E, 0x47}) {
            return "image/png"
        }
    }
    
    // Por defecto
    return "application/octet-stream"
}
