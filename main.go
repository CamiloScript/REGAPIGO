package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"github.com/CamiloScript/REGAPIGO/src/APIGuest/APIController"
	"github.com/CamiloScript/REGAPIGO/bd/MongoDB"
)

func main() {
	// Inicializamos el servidor de Gin
	router := gin.Default()

	// Conexión a la base de datos
	_, err := MongoDB.ConexionDB()
	if err != nil {
		log.Fatal("Error al conectar a la base de datos: ", err)
	}
	defer MongoDB.CerrarConexion() // Cerramos la conexión cuando termine el programa

	// Definimos las rutas para las peticiones HTTP
	router.GET("/guests", APIController.GetGuests)
	router.GET("/guests/:id", APIController.GetGuestByID)
	router.POST("/guests", APIController.CreateGuest)
	router.PUT("/guests/:id", APIController.UpdateGuest)
	router.DELETE("/guests/:id", APIController.DeleteGuest)

	// Iniciamos el servidor
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Error al iniciar el servidor: ", err)
	}

	// Parar el servidor
	// Mongodb.CerrarConexion()
}
