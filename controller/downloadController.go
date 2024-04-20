package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"github.com/gorilla/mux"
)

// DownloadFileHandler descarga un archivo por su ID
func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID := vars["fileId"]
	userID := vars["userId"]
	filePath := fmt.Sprintf("/mnt/nfs/%s/%s", userID, fileID)

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Obtener el tipo de contenido del archivo basado en su contenido
	contentType := http.DetectContentType(nil)

	// Establecer el tipo de contenido adecuado
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Set appropriate headers
	w.Header().Set("Content-Disposition", "attachment; filename="+fileID)
	w.Header().Set("Content-Type", contentType)

	// Copiar el archivo al escritor de respuesta
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
