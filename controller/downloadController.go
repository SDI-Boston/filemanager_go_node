package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"path/filepath"
	"mime"
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

    // Obtener la extensión del archivo
    fileExt := filepath.Ext(fileID)

    // Establecer el tipo de contenido basado en la extensión del archivo
    contentType := mime.TypeByExtension(fileExt)
    if contentType == "" {
        contentType = "application/octet-stream"
    }

    // Establecer los encabezados de respuesta
    w.Header().Set("Content-Disposition", "attachment; filename="+fileID)
    w.Header().Set("Content-Type", contentType)

    // Copiar el archivo al escritor de respuesta
    _, err = io.Copy(w, file)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

