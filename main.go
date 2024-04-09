package main

import (
	"fmt"

	"github.com/SDI-Boston/filemanager_go_node/controller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	// Configuración del servidor Echo para las posteriores rutas
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Rutas del servidor

	// Test de conexión con el servidor NFS
	if err := controller.TestNFSConnection(); err != nil {
		// Manejar el error aquí, por ejemplo, imprimirlo y salir
		fmt.Println("Error:", err)
		return
	}

	// Inicia el servidor Echo
	//e.Logger.Fatal(e.Start(":5000"))
}
