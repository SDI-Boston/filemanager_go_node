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

	// Intenta descargar desde la ruta principal primero
	filePath := fmt.Sprintf("/mnt/nfs/%s/%s", userID, fileID)
	file, err := downloadFile(filePath, w)
	if err != nil {
		// Si no se pudo descargar desde la ruta principal, intenta descargar desde la ruta de respaldo
		backupFilePath := fmt.Sprintf("/mnt/nfs_backup/%s/%s", userID, fileID)
		file, err = downloadFile(backupFilePath, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	defer file.Close()

	// Set appropriate headers
	w.Header().Set("Content-Disposition", "attachment; filename="+fileID)

	// Copy the file to the response writer
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func downloadFile(filePath string, w http.ResponseWriter) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}
