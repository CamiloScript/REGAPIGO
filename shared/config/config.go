package config

import (
    "encoding/json"
    "os"
    "strconv"
    "log"
)

// Config es una estructura que centraliza toda la configuración de la aplicación.
type Config struct {
    Puerto            string  // Puerto en el que se ejecutará el servidor
    AlfrescoBaseURL   string  // URL base para la conexión con Alfresco
    AlfrescoAPIKey    string  // API Key para autenticación con Alfresco
    MaxFileSize       int64   // Tamaño máximo permitido para archivos (en bytes)
    LogLevel          string  // Nivel de logging (INFO, DEBUG, ERROR, etc.)
    SeparadorLog      string  // Carácter separador para los logs
    ClaveSesion       string  // Clave para firmar las cookies de sesión
    TiempoSesion      int     // Duración de la sesión en horas
    MongoURI          string  // URL de conexión a MongoDB
    MongoDatabase     string  // Nombre de la base de datos en MongoDB
    MongoCollection   string  // Nombre de la colección en MongoDB
    AuthUser          string  // Usuario de autenticación interna
    AuthPassword      string  // Contraseña de autenticación interna
    ApiKey            string  // API Key para solicitudes externas
}

// CargarConfiguracion carga la configuración desde el archivo appsettings.json.
func CargarConfiguracion() *Config {
    // Abrir el archivo appsettings.json
    file, err := os.Open("./settings/appsettings.json")
    if err != nil {
        log.Fatalf("No se pudo abrir el archivo appsettings.json: %v", err)
    }
    defer file.Close()

    // Leer el contenido del archivo
    var configMap map[string]string
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&configMap); err != nil {
        log.Fatalf("Error al decodificar appsettings.json: %v", err)
    }

    // Retornar una nueva instancia de Config con valores cargados
    return &Config{
        Puerto:            configMap["PORT"],
        AlfrescoBaseURL:   configMap["ALFRESCO_BASE_URL"],
        AlfrescoAPIKey:    configMap["ALFRESCO_API_KEY"],
        MaxFileSize:       getEnvAsInt64(configMap["MAX_FILE_SIZE"], 10*1024*1024),
        LogLevel:          configMap["LOG_LEVEL"],
        SeparadorLog:      configMap["LOG_SEPARATOR"],
        ClaveSesion:       configMap["SESSION_KEY"],
        TiempoSesion:      getEnvAsInt(configMap["SESSION_DURATION"], 24),
        MongoURI:          configMap["MONGODB_URI"],
        MongoDatabase:     configMap["MONGODB_DATABASE"],
        MongoCollection:   configMap["MONGODB_COLLECTION"],
        AuthUser:          configMap["AUTH_USER"],
        AuthPassword:      configMap["AUTH_PASSWORD"],
        ApiKey:            configMap["API_KEY"],
    }
}

// getEnvAsInt64 es una función auxiliar que convierte un string a int64.
func getEnvAsInt64(value string, defaultValue int64) int64 {
    if value == "" {
        return defaultValue
    }
    intValue, err := strconv.ParseInt(value, 10, 64)
    if err != nil {
        return defaultValue
    }
    return intValue
}

// getEnvAsInt es una función auxiliar que convierte un string a int.
func getEnvAsInt(value string, defaultValue int) int {
    if value == "" {
        return defaultValue
    }
    intValue, err := strconv.Atoi(value)
    if err != nil {
        return defaultValue
    }
    return intValue
}