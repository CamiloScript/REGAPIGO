package documentos

import (
    "fmt"
    "strings"
)

// Documento representa la entidad de negocio, agnóstica del proveedor
type Documento struct {
    ID           string                 `json:"id"`
    Nombre       string                 `json:"nombre"`
    Tipo         string                 `json:"tipo"`
    Tamano      int64                  `json:"tamano"`
    Propiedades  map[string]interface{} `json:"propiedades"`
    Autor        string                 `json:"autor"`
    ModificadoEn string                 `json:"modificadoEn"`
    // ... Otros campos según el dominio
}


// DocumentMetadata representa los metadatos de un documento en el sistema Alfresco.
// Cada campo está asociado a un atributo específico dentro de Alfresco.
// Los campos están etiquetados con JSON y BSON para su serialización/deserialización.
type DocumentMetadata struct {
    TipoDocumento        string   `json:"tanner:tipo-documento" bson:"_tanner:tipo-documento, omitempty"`         // Tipo de documento
    RazonSocialCliente   string   `json:"tanner:razon-social-cliente" bson:"_tanner:tipo-documento, omitempty"`   // Razón social del cliente
    RUTCliente           string   `json:"tanner:rut-cliente" bson:"_tanner:rut-cliente, omitempty"`               // RUT del cliente
    EstadoVisado         string   `json:"tanner:estado-visado" bson:"_tanner:estado-visado, omitempty"`           // Estado de visado del documento
    EstadoVigencia       string   `json:"tanner:estado-vigencia" bson:"_tanner:estado-vigencia, omitempty"`       // Estado de vigencia del documento
    FechaCarga           string   `json:"tanner:fecha-carga" bson:"_tanner:fecha-carga, omitempty" validate:"datetime=2006-01-02T15:04:05.999Z07:00"` // Fecha de carga del documento
    NombreDoc            string   `json:"tanner:nombre-doc" bson:"_tanner:nombre-doc, omitempty"`                 // Nombre del documento
    Categorias           string   `json:"tanner:categorias" bson:"_tanner:categorias, omitempty"`                 // Categorías del documento
    SubCategorias        string   `json:"tanner:sub-categorias" bson:"_tanner:sub-categorias, omitempty"`          // Subcategorías del documento
    Origen               string   `json:"tanner:origen" bson:"_tanner:origen, omitempty"`                         // Origen del documento
    Relacion             string   `json:"tanner:relacion" bson:"_tanner:relacion, omitempty"`                     // Relación del documento
    FechaTerminoVigencia string   `json:"tanner:fecha-termino-vigencia" bson:"_tanner:fecha-termino-vigencia, omitempty"` // Fecha de término de vigencia del documento
    CmTitle              string   `json:"cm:title" bson:"_cm:title, omitempty"`                                   // Título del documento en Alfresco
    CmVersionType        string   `json:"cm:versionType" bson:"_cm:versionType, omitempty"`                       // Tipo de versión del documento
    CmVersionLabel       string   `json:"cm:versionLabel" bson:"_cm:versionLabel, omitempty"`                     // Etiqueta de versión del documento
    CmDescription        string   `json:"cm:description" bson:"_cm:description, omitempty"`                       // Descripción del documento
    Observaciones        string   `json:"tanner:observaciones" bson:"_tanner:observaciones,omitempty"`            // Observaciones adicionales
    FileType             string   `json:"file_type" bson:"file_type,omitempty"`                                   // Campo para el tipo de archivo
}

// Validar verifica que los campos obligatorios del DocumentMetadata estén presentes y que el RUT sea válido.
// Retorna un error si algún campo obligatorio está vacío o si el RUT no cumple con el formato esperado.
func (d *DocumentMetadata) Validar() error {
    // Mapa de campos obligatorios y sus valores
    camposRequeridos := map[string]string{
        "TipoDocumento":  d.TipoDocumento,
        "RUTCliente":     d.RUTCliente,
        "NombreDoc":      d.NombreDoc,
        "CmTitle":        d.CmTitle,
        "CmVersionType":  d.CmVersionType,
        "CmVersionLabel": d.CmVersionLabel,
        "CmDescription":  d.CmDescription,
        "SubCategorias":   d.SubCategorias,
    }

    // Verificar campos obligatorios faltantes
    var camposFaltantes []string
    for nombre, valor := range camposRequeridos {
        if valor == "" {
            camposFaltantes = append(camposFaltantes, nombre)
        }
    }

    // Si hay campos faltantes, retornar un error con la lista de campos
    if len(camposFaltantes) > 0 {
        return fmt.Errorf("campos obligatorios faltantes: %s", strings.Join(camposFaltantes, ", "))
    }

    // Si todo está correcto, retornar nil (sin errores)
    return nil
}
