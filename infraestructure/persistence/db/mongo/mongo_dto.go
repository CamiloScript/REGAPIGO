package mongo

// DocumentoMongoDTO define la estructura personalizada para almacenar documentos en MongoDB.
// Esta estructura mapea los datos del documento a los campos correspondientes en la base de datos.
type DocumentoMongoDTO struct {
    ID            string                 `bson:"id_archivo"`      // ID del documento en Alfresco
    NombreArchivo string                 `bson:"nombre_archivo"`  // Nombre del archivo
    TipoArchivo   string                 `bson:"tipo_archivo"`    // Tipo de archivo, por ejemplo: "application/pdf"
    FechaCarga    string                 `bson:"fecha_carga"`     // Fecha en que se cargó el archivo
    Metadatos     map[string]interface{} `bson:"metadatos"`       // Metadatos adicionales del documento
}