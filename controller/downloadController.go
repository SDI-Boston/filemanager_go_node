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

	// Intentar abrir el archivo desde /mnt/nfs
	filePath := fmt.Sprintf("/mnt/nfs/%s/%s", userID, fileID)
	file, err := os.Open(filePath)
	if err != nil {
		// Si falla, intentar abrir el archivo desde /mnt/nfs_backup
		filePath = fmt.Sprintf("/mnt/nfs_backup/%s/%s", userID, fileID)
		file, err = os.Open(filePath)
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
