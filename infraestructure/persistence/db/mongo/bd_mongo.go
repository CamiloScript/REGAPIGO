package mongo

import (
    "context"
    "fmt"
    "time"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "github.com/CamiloScript/REGAPIGO/shared/config"
    "encoding/json"
    "github.com/CamiloScript/REGAPIGO/shared/logger"
    "github.com/CamiloScript/REGAPIGO/application/documento"
)

// Variables globales para reutilizar la conexión y configuración
var (
    client *mongo.Client
    cfg    = config.CargarConfiguracion() // Cargar configuración desde archivo
)

// ConexionDB establece o reutiliza una conexión a MongoDB.
// Retorna un cliente funcional o un error en caso de fallo.
func ConexionDB() (*mongo.Client, error) {
    // Reutilizar conexión existente si ya está establecida
    if client != nil {
        return client, nil
    }

    // Configurar opciones del cliente usando la URI de MongoDB desde la configuración
    opts := options.Client().ApplyURI(cfg.MongoURI)
    opts.SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))

    // Crear un contexto con timeout para evitar bloqueos
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Establecer conexión con MongoDB
    var err error
    client, err = mongo.Connect(ctx, opts)
    if err != nil {
        return nil, fmt.Errorf("fallo al conectar: %v", err)
    }

    // Verificar que la conexión esté activa
    if err := client.Ping(ctx, nil); err != nil {
        return nil, fmt.Errorf("fallo al verificar conexión: %v", err)
    }

    fmt.Println("✅ Conexión a MongoDB establecida")
    return client, nil
}

// InsertarDocumento inserta un documento en la colección especificada.
// Parámetros:
//   - document: Datos a insertar (debe ser serializable a BSON).
// Retorna un error si falla la operación.
func InsertarDocumento(document interface{}) error {
    // Obtener conexión a MongoDB
    client, err := ConexionDB()
    if err != nil {
        return fmt.Errorf("error obteniendo cliente: %v", err)
    }

    // Crear un contexto con timeout para la operación de inserción
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Seleccionar base de datos y colección desde la configuración
    collection := client.Database(cfg.MongoDatabase).Collection(cfg.MongoCollection)

    // Ejecutar la inserción del documento
    _, err = collection.InsertOne(ctx, document)
    return err
}

// GuardarEnMongoDB persiste un documento en MongoDB.
// Parámetros:
//   - respuesta: Mapa con los datos del documento a guardar.
//   - log: Logger para registrar eventos y errores.
// Retorna un error si falla la operación.
func GuardarEnMongoDB(respuesta map[string]interface{}, log *logger.Registrador) error {
    // 1. Mapear respuesta al DTO de Alfresco usando JSON
    var alfrescoDoc documento.AlfrescoDocumentDTO

    // Convertir la respuesta del servicio a JSON
    respuestaJSON, err := json.Marshal(respuesta)
    if err != nil {
        log.Error("Error serializando respuesta", map[string]interface{}{
            "error": err.Error(),
            "stack": "paso1",
        })
        return err
    }

    // Deserializar el JSON a la estructura AlfrescoDocumentDTO
    if err := json.Unmarshal(respuestaJSON, &alfrescoDoc); err != nil {
        log.Error("Error mapeando a DTO", map[string]interface{}{
            "error": err.Error(),
            "stack": "paso2",
        })
        return err
    }

    // 2. Validar campos críticos antes de guardar
    if alfrescoDoc.Entry.ID == "" {
        log.Warn("Documento sin ID - No se guardará en MongoDB",
            map[string]interface{}{"origen": "Alfresco", "stack": "validación"})
        return fmt.Errorf("documento sin ID válido")
    }

    // 3. Convertir alfrescoDoc.Entry al DTO de MongoDB
    mongoDoc := DocumentoMongoDTO{
        ID:            alfrescoDoc.Entry.ID,
        NombreArchivo: alfrescoDoc.Entry.Name,
        TipoArchivo:   alfrescoDoc.Entry.Content.MimeType,
        FechaCarga:    alfrescoDoc.Entry.Properties["tanner:fecha-carga"].(string),
        Metadatos: map[string]interface{}{
            "nombre_documento":      alfrescoDoc.Entry.Properties["tanner:nombre-doc"],
            "tipo_documento":        alfrescoDoc.Entry.Properties["tanner:tipo-documento"],
            "razon_social_cliente":  alfrescoDoc.Entry.Properties["tanner:razon-social-cliente"],
            "rut_cliente":           alfrescoDoc.Entry.Properties["tanner:rut-cliente"],
            "estado_vigencia":       alfrescoDoc.Entry.Properties["tanner:estado-vigencia"],
            "fecha_carga":           alfrescoDoc.Entry.Properties["tanner:fecha-carga"],
            // Otros campos según necesidad
        },
    }

    // 4. Insertar el documento en MongoDB
    if err := InsertarDocumento(mongoDoc); err != nil {
        log.Error("Fallo al guardar en MongoDB",
            map[string]interface{}{
                "error":     err.Error(),
                "doc_id":    alfrescoDoc.Entry.ID,
            })
        return err
    }

    // 5. Registrar éxito de la operación
    log.Info("Documento persistido exitosamente",
        map[string]interface{}{
            "id":        mongoDoc.ID,
            "tamaño_mb": alfrescoDoc.Entry.Content.SizeInBytes / (1024 * 1024),
        })

    return nil
}

// BuscarDocumento busca un documento en MongoDB usando un filtro.
// Parámetros:
//   - filtro: Mapa con criterios de búsqueda (ej: {"metadatos.rut_cliente": "12345678-9"}).
// Retorna:
//   - idFile: ID del documento encontrado en Alfresco.
//   - error: Mensaje en caso de fallo.
func BuscarDocumento(filtro map[string]interface{}) (string, error) {
    // 1. Obtener conexión a MongoDB
    client, err := ConexionDB()
    if err != nil {
        return "", fmt.Errorf("error de conexión: %v", err)
    }

    // 2. Seleccionar la colección desde la configuración
    collection := client.Database(cfg.MongoDatabase).Collection(cfg.MongoCollection)

    // 3. Configurar contexto con timeout para la operación de búsqueda
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // 4. Ejecutar la consulta usando el filtro proporcionado
    var resultado DocumentoMongoDTO
    err = collection.FindOne(ctx, filtro).Decode(&resultado)
    if err != nil {
        return "", fmt.Errorf("documento no encontrado: %v", err)
    }

    // 5. Validar que el ID del documento sea válido
    if resultado.ID == "" {
        return "", fmt.Errorf("ID de archivo inválido en MongoDB")
    }

    return resultado.ID, nil
}