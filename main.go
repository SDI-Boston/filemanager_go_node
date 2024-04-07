package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	//"tu_paquete/grpc" // Importa la lógica del servidor gRPC
)

func main() {
	// Configuración del servidor Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Rutas del servidor
	//e.POST("/upload", grpc.UploadHandler) // Test de la ruta y su metodo

	// Inicia el servidor Echo
	e.Logger.Fatal(e.Start(":5000"))
}
