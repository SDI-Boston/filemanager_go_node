package main

import (
	"fmt"

	"github.com/SDI-Boston/filemanager_go_node/controller"
)

func main() {

	// Test de conexión y creación con el servidor NFS
	if err := controller.TestNFSConnection(); err != nil {
		// Manejar el error aquí, por ejemplo, imprimirlo y salir
		fmt.Println("Error:", err)
		return
	}

	/*go func() {
		if err := http.ListenAndServeTLS(":5000", "cert.pem", "key.pem", router); err != nil {
			fmt.Printf("Error al iniciar el servidor HTTPS: %s\n", err)
		}
	}()

	fmt.Println("Servidor en ejecución en https://localhost:5000")*/

	//select {}

}
