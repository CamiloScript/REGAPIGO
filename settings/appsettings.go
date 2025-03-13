package settings

import (
    "encoding/json"
    "io"
    "log"
    "os"
)

// Settings representa la configuración completa
type Settings struct {
    Port            string `json:"PORT"`
    AlfrescoBaseURL string `json:"ALFRESCO_BASE_URL"`
    AlfrescoAPIKey  string `json:"ALFRESCO_API_KEY"`
    LogSeparator    string `json:"LOG_SEPARATOR"`
    LogLevel        string `json:"LOG_LEVEL"`
    MaxFileSize     string `json:"MAX_FILE_SIZE"`
    MongoURI        string `json:"MONGODB_URI"`
    MongoDatabase   string `json:"MONGODB_DATABASE"`
    MongoCollection string `json:"MONGODB_COLLECTION"`
    AuthUser        string `json:"AUTH_USER"`
    AuthPassword    string `json:"AUTH_PASSWORD"`
    ApiKey          string `json:"API_KEY"`
}

// LoadConfiguration carga la configuración desde un archivo JSON
func LoadConfiguration(route string) *Settings {
    settings := &Settings{}

    // Abrir el archivo JSON
    configFile, err := os.Open(route)
    if err != nil {
        log.Fatalf("No se pudo abrir el archivo de configuración: %v", err)
    }
    defer configFile.Close()

    // Leer el contenido del archivo
    bytes, err := io.ReadAll(configFile)
    if err != nil {
        log.Fatalf("Error al leer el archivo de configuración: %v", err)
    }

    // Decodificar el JSON en la estructura Settings
    if err := json.Unmarshal(bytes, settings); err != nil {
        log.Fatalf("Error al decodificar el archivo de configuración: %v", err)
    }

    return settings
}