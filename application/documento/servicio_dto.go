package documento

// AlfrescoDocumentDTO representa la respuesta de subida y detalles de un documento en Alfresco.
type AlfrescoDocumentDTO struct {
    Entry struct {
        ID           string                 `json:"id"`           // ID único del documento
        Name         string                 `json:"name"`         // Nombre del documento
        NodeType     string                 `json:"nodeType"`     // Tipo de nodo (ej. "cm:content")
        IsFolder     bool                   `json:"isFolder"`     // Indica si es una carpeta
        IsFile       bool                   `json:"isFile"`       // Indica si es un archivo
        IsLocked     interface{}            `json:"isLocked"`     // Indica si el documento está bloqueado (puede ser null)
        ModifiedAt   string                 `json:"modifiedAt"`   // Fecha de última modificación
        ModifiedByUser struct {
            ID          string `json:"id"`          // ID del usuario que modificó el documento
            DisplayName string `json:"displayName"` // Nombre del usuario que modificó el documento
        } `json:"modifiedByUser"` // Información del usuario que modificó el documento
        CreatedAt    string                 `json:"createdAt"`    // Fecha de creación del documento
        CreatedByUser struct {
            ID          string `json:"id"`          // ID del usuario que creó el documento
            DisplayName string `json:"displayName"` // Nombre del usuario que creó el documento
        } `json:"createdByUser"` // Información del usuario que creó el documento
        ParentId     string                 `json:"parentId"`     // ID del nodo padre (carpeta contenedora)
        Content      struct {
            MimeType     string `json:"mimeType"`     // Tipo MIME del archivo (ej. "application/pdf")
            MimeTypeName string `json:"mimeTypeName,omitempty"` // Nombre del tipo MIME (opcional)
            SizeInBytes  int64  `json:"sizeInBytes"`  // Tamaño del archivo en bytes
            Encoding     string `json:"encoding,omitempty"` // Codificación del archivo (opcional)
        } `json:"content"` // Información sobre el contenido del archivo
        Properties   map[string]interface{} `json:"properties"`   // Propiedades adicionales del documento
        AspectNames  []string               `json:"aspectNames,omitempty"` // Nombres de aspectos aplicados (opcional)
        AllowableOperations []string        `json:"allowableOperations,omitempty"` // Operaciones permitidas (opcional)
        Path         interface{}            `json:"path"` // Ruta del documento (puede ser null o un objeto)
    } `json:"entry"` // Entrada principal que contiene los detalles del documento
}

// AlfrescoDocumentListDTO representa la respuesta del listado de documentos en Alfresco.
type AlfrescoDocumentListDTO []struct {
    Entry struct {
        ID           string                 `json:"id"`           // ID único del documento
        Name         string                 `json:"name"`         // Nombre del documento
        NodeType     string                 `json:"nodeType"`     // Tipo de nodo (ej. "cm:content")
        IsFolder     bool                   `json:"isFolder"`     // Indica si es una carpeta
        IsFile       bool                   `json:"isFile"`       // Indica si es un archivo
        ModifiedAt   string                 `json:"modifiedAt"`   // Fecha de última modificación
        CreatedAt    string                 `json:"createdAt"`    // Fecha de creación del documento
        ParentId     string                 `json:"parentId"`     // ID del nodo padre (carpeta contenedora)
        Content      struct {
            MimeType     string `json:"mimeType"`     // Tipo MIME del archivo (ej. "application/pdf")
            SizeInBytes  int64  `json:"sizeInBytes"`  // Tamaño del archivo en bytes
        } `json:"content"` // Información sobre el contenido del archivo
        Properties   map[string]interface{} `json:"properties"`   // Propiedades adicionales del documento
    } `json:"entry"` // Entrada principal que contiene los detalles del documento
}

// AlfrescoDownloadResponse representa la respuesta de descarga de un documento en Alfresco.
type AlfrescoDownloadResponse struct {
    Content        []byte   `json:"-"` // Contenido binario del archivo (no se serializa en JSON)
    FileName       string   `json:"fileName"` // Nombre del archivo obtenido del header Content-Disposition
    MimeType       string   `json:"mimeType"` // Tipo MIME del archivo obtenido del header
}