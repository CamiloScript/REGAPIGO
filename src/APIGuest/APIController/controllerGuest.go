package APIController

import (
	"fmt"
	"net/http"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/CamiloScript/REGAPIGO/src/APIGuest/APIStruct"
	"github.com/CamiloScript/REGAPIGO/bd/MongoDB"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

)

// GetGuests maneja la solicitud GET para obtener todos los huéspedes
func GetGuests(c *gin.Context) {
	
	// Consultamos la base de datos
	collection := MongoDB.Cliente.Database("APIREGDB").Collection("Guests")
	cursor, err := collection.Find(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al obtener los huéspedes", nil))
		return
	}
	defer cursor.Close(c)

	var guests []APIStruct.Guest
	if err := cursor.All(c, &guests); err != nil {
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al procesar los huéspedes", nil))
		return
	}

	c.JSON(http.StatusOK, FormatResponse("success", "Lista de huéspedes obtenida correctamente", guests))
}

// GetGuestByID maneja la solicitud GET para obtener un huésped por su ID
func GetGuestByID(c *gin.Context) {
    id := c.Param("id")
    fmt.Println("ID recibido:", id)  // Agregamos una línea de depuración

    // Convertimos el id de string a ObjectID
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        fmt.Println("Error al convertir el ID a ObjectID:", err)  // Depuración
        c.JSON(http.StatusBadRequest, FormatResponse("error", "ID inválido", nil))
        return
    }

    fmt.Println("ObjectID convertido:", objectID)  // Depuración

    // Consultamos la base de datos
    collection := MongoDB.Cliente.Database("APIREGDB").Collection("Guests")
    var guest APIStruct.Guest
    err = collection.FindOne(c, bson.M{"_id": objectID}).Decode(&guest)

    if err == mongo.ErrNoDocuments {
        fmt.Println("No se encontró el huésped con ese ID.")  // Depuración
        c.JSON(http.StatusNotFound, FormatResponse("error", "Huésped no encontrado", nil))
        return
    } else if err != nil {
        fmt.Println("Error al obtener el huésped:", err)  // Depuración
        c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al obtener el huésped", nil))
        return
    }

    fmt.Println("Huésped encontrado:", guest)  // Depuración

    c.JSON(http.StatusOK, FormatResponse("success", "Huésped encontrado", guest))
}


// CreateGuest maneja la solicitud POST para crear un nuevo huésped
func CreateGuest(c *gin.Context) {
	var newGuest APIStruct.Guest
	if err := c.ShouldBindJSON(&newGuest); err != nil {
		c.JSON(http.StatusBadRequest, FormatResponse("error", "Datos inválidos", nil))
		return
	}

	// Insertamos el nuevo huésped en la base de datos
	collection := MongoDB.Cliente.Database("APIREGDB").Collection("Guests")
	_, err := collection.InsertOne(c, newGuest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al crear el huésped", nil))
		return
	}

	c.JSON(http.StatusCreated, FormatResponse("success", "Huésped creado correctamente", newGuest))
}

// UpdateGuest maneja la solicitud PUT para actualizar los datos de un huésped
func UpdateGuest(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResponse("error", "ID inválido", nil))
		return
	}

	var updatedGuest APIStruct.Guest
	if err := c.ShouldBindJSON(&updatedGuest); err != nil {
		c.JSON(http.StatusBadRequest, FormatResponse("error", "Datos inválidos", nil))
		return
	}

	// Actualizamos el huésped en la base de datos
	collection := MongoDB.Cliente.Database("APIREGDB").Collection("Guests")
	result, err := collection.UpdateOne(c, bson.M{"_id": objectID}, bson.M{"$set": updatedGuest})

	if err != nil {
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al actualizar el huésped", nil))
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, FormatResponse("error", "Huésped no encontrado", nil))
		return
	}

	c.JSON(http.StatusOK, FormatResponse("success", "Huésped actualizado correctamente", updatedGuest))
}

// DeleteGuest maneja la solicitud DELETE para eliminar un huésped por su ID
func DeleteGuest(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, FormatResponse("error", "ID inválido", nil))
		return
	}

	// Eliminamos el huésped de la base de datos
	collection := MongoDB.Cliente.Database("APIREGDB").Collection("Guests")
	result, err := collection.DeleteOne(c, bson.M{"_id": objectID})

	if err != nil {
		c.JSON(http.StatusInternalServerError, FormatResponse("error", "Error al eliminar el huésped", nil))
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, FormatResponse("error", "Huésped no encontrado", nil))
		return
	}

	c.JSON(http.StatusOK, FormatResponse("success", "Huésped eliminado correctamente", nil))
}


// FormatResponse les aplica el formato presente en APIStruct.ApiResponse a las respuestas de la API
//totalmente perzonalizable, y englobado en una función movil
func FormatResponse(status, message string, data interface{}) APIStruct.ApiResponse {
    return APIStruct.ApiResponse{
        Status:  status,
        Message: message,
        Data:    data,
    }
}
