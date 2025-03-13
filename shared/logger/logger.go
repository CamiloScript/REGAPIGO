package logger

import (
    "os"
    "strings"
    "github.com/rs/zerolog"
    "fmt"
)

// Registrador estructura principal para el registro personalizado
type Registrador struct {
    zerolog.Logger    // Embebe el logger de zerolog para utilizar sus funcionalidades
    nombreServicio string // Nombre del servicio asociado al registrador
}

// NuevoRegistrador crea una nueva instancia del registrador con formato personalizado
// Parámetros:
//   - nombreServicio: Nombre del servicio para identificar los logs.
//   - separador: Carácter separador para los valores de los campos en los logs.
// Retorna una instancia configurada de Registrador.
func NuevoRegistrador(nombreServicio, separador string) *Registrador {
    // Define la salida del logger con un formato de consola
    salida := zerolog.ConsoleWriter{
        Out:        os.Stdout,                      // La salida será la consola estándar
        TimeFormat: "2006-01-02T15:04:05Z07:00",   // Formato de fecha y hora
    }

    // Personaliza la forma en que se muestran los valores de los campos
    salida.FormatFieldValue = func(i interface{}) string {
        return separador + i.(string) // Agrega un separador al valor del campo
    }

    // Inicializa el logger con el servicio asociado y la marca de tiempo
    registrador := zerolog.New(salida).
        With().
        Timestamp(). // Agrega la marca de tiempo automáticamente
        Str("servicio", nombreServicio). // Agrega el nombre del servicio a los logs
        Logger()

    return &Registrador{
        Logger:          registrador,
        nombreServicio: nombreServicio,
    }
}

// RegistrarDebug registra un mensaje de nivel DEBUG
// Parámetros:
//   - mensaje: Mensaje a registrar.
//   - contexto: Contexto adicional en forma de mapa.
func (r *Registrador) Debug(mensaje string, contexto map[string]interface{}) {
    // Convertir todos los valores a string para evitar problemas de tipo (por ejemplo, json.Number)
    contextoString := make(map[string]interface{})
    for k, v := range contexto {
        contextoString[k] = fmt.Sprintf("%v", v)
    }
    r.Logger.Debug().Fields(contextoString).Msg(mensaje) // Registra un mensaje con contexto adicional
}

// RegistrarInfo registra un mensaje de nivel INFO
// Parámetros:
//   - mensaje: Mensaje a registrar.
//   - contexto: Contexto adicional en forma de mapa.
func (r *Registrador) Info(mensaje string, contexto map[string]interface{}) {
    // Convertir todos los valores a string para evitar problemas de tipo (por ejemplo, json.Number)
    contextoString := make(map[string]interface{})
    for k, v := range contexto {
        contextoString[k] = fmt.Sprintf("%v", v)
    }
    r.Logger.Info().Fields(contextoString).Msg(mensaje) // Registra un mensaje informativo
}

// RegistrarAdvertencia registra un mensaje de nivel WARN
// Parámetros:
//   - mensaje: Mensaje a registrar.
//   - contexto: Contexto adicional en forma de mapa.
func (r *Registrador) Warn(mensaje string, contexto map[string]interface{}) {
    // Convertir todos los valores a string para evitar problemas de tipo (por ejemplo, json.Number)
    contextoString := make(map[string]interface{})
    for k, v := range contexto {
        contextoString[k] = fmt.Sprintf("%v", v)
    }
    r.Logger.Warn().Fields(contextoString).Msg(mensaje) // Registra un mensaje de advertencia
}

// RegistrarError registra un mensaje de nivel ERROR
// Parámetros:
//   - mensaje: Mensaje a registrar.
//   - contexto: Contexto adicional en forma de mapa.
func (r *Registrador) Error(mensaje string, contexto map[string]interface{}) {
    // Convertir todos los valores a string para evitar problemas de tipo (por ejemplo, json.Number)
    contextoString := make(map[string]interface{})
    for k, v := range contexto {
        contextoString[k] = fmt.Sprintf("%v", v)
    }
    r.Logger.Error().Fields(contextoString).Msg(mensaje) // Registra un mensaje de error
}

// RegistrarFatal registra un mensaje de nivel FATAL y termina la aplicación
// Parámetros:
//   - mensaje: Mensaje a registrar.
//   - contexto: Contexto adicional en forma de mapa.
func (r *Registrador) Fatal(mensaje string, contexto map[string]interface{}) {
    // Convertir todos los valores a string para evitar problemas de tipo (por ejemplo, json.Number)
    contextoString := make(map[string]interface{})
    for k, v := range contexto {
        contextoString[k] = fmt.Sprintf("%v", v)
    }
    r.Logger.Fatal().Fields(contextoString).Msg(mensaje) // Registra un mensaje fatal y detiene la ejecución
}

// RegistrarPanico registra un mensaje de nivel PANIC y provoca un pánico en la aplicación
// Parámetros:
//   - mensaje: Mensaje a registrar.
//   - contexto: Contexto adicional en forma de mapa.
func (r *Registrador) Panic(mensaje string, contexto map[string]interface{}) {
    // Convertir todos los valores a string para evitar problemas de tipo (por ejemplo, json.Number)
    contextoString := make(map[string]interface{})
    for k, v := range contexto {
        contextoString[k] = fmt.Sprintf("%v", v)
    }
    r.Logger.Panic().Fields(contextoString).Msg(mensaje) // Registra un mensaje de pánico
}

// EstablecerNivel configura el nivel de registro global
// Parámetros:
//   - nivel: Nivel de registro a establecer (DEBUG, INFO, WARN, ERROR, FATAL, PANIC).
func (r *Registrador) EstablecerNivel(nivel string) {
    switch strings.ToUpper(nivel) { // Convierte la entrada a mayúsculas para evitar errores
    case "DEBUG":
        zerolog.SetGlobalLevel(zerolog.DebugLevel) // Nivel de depuración
    case "INFO":
        zerolog.SetGlobalLevel(zerolog.InfoLevel) // Nivel informativo
    case "WARN":
        zerolog.SetGlobalLevel(zerolog.WarnLevel) // Nivel de advertencia
    case "ERROR":
        zerolog.SetGlobalLevel(zerolog.ErrorLevel) // Nivel de error
    case "FATAL":
        zerolog.SetGlobalLevel(zerolog.FatalLevel) // Nivel fatal
    case "PANIC":
        zerolog.SetGlobalLevel(zerolog.PanicLevel) // Nivel pánico
    default:
        zerolog.SetGlobalLevel(zerolog.InfoLevel)  // Nivel por defecto (INFO)
    }
}

// EstablecerFormatoJSON cambia el formato del registro a JSON
func (r *Registrador) EstablecerFormatoJSON() {
    r.Logger = zerolog.New(os.Stdout). // Define la salida como la consola estándar
        With().
        Timestamp(). // Agrega la marca de tiempo automáticamente
        Str("servicio", r.nombreServicio). // Incluye el nombre del servicio en los logs
        Logger()
}