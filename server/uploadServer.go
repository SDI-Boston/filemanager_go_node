package server

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	pb "github.com/SDI-Boston/filemanager_go_node/proto"
)

type FileService struct {
	pb.UnimplementedFileServiceServer
}

func (s *FileService) Upload(stream pb.FileService_UploadServer) error {
	// Leer el request desde el flujo de datos
	req, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("failed to receive upload request: %w", err)
	}

	// Validar hash
	if err := validateHash(req); err != nil {
		return fmt.Errorf("hash validation failed: %w", err)
	}

	// Subir archivo
	filePath, err := uploadToNFS(req)
	if err != nil {
		return fmt.Errorf("failed to upload file to NFS: %w", err)
	}

	// Respuesta
	err = stream.SendAndClose(&pb.FileUploadResponse{
		FileId: req.FileId,
		Urls:   []string{filePath},
	})
	if err != nil {
		return fmt.Errorf("failed to send upload response: %w", err)
	}

	return nil
}

func validateHash(req *pb.FileUploadRequest) error {
	// Calcular el hash SHA256 del contenido del archivo
	hash := sha256.New()
	hash.Write(req.BinaryFile)
	calculatedHash := fmt.Sprintf("%x", hash.Sum(nil))

	fmt.Println("Received Hash:", req.FileHash)
	fmt.Println("Calculated Hash:", calculatedHash)

	// Comparar el hash calculado con el hash proporcionado en la solicitud
	if calculatedHash != req.FileHash {
		return fmt.Errorf("file hash mismatch")
	}

	return nil
}


func uploadToNFS(req *pb.FileUploadRequest) (string, error) {
	// Si el usuario nunca ha creado un archivo, crear un directorio para el usuario
	userPath := fmt.Sprintf("/mnt/nfs/%s", req.OwnerId)
	if _, err := os.Stat(userPath); os.IsNotExist(err) {
		err := os.Mkdir(userPath, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create user directory: %w", err)
		}
	}

	// Extraer la extensión del nombre del archivo
	fileExtension := filepath.Ext(req.FileName)

	fileName := req.FileId + fileExtension
	filePath := filepath.Join(userPath, fileName)

	fileUpload, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	defer fileUpload.Close()

	// Decodificar el contenido del binario que viene en base64
	decodedContent, err := base64.StdEncoding.DecodeString(string(req.BinaryFile))
	if err != nil {
		return "", fmt.Errorf("failed to decode binary content: %w", err)
	}

	// Escribir el contenido decodificado en el archivo
	_, err = fileUpload.Write(decodedContent)
	if err != nil {
		return "", fmt.Errorf("failed to write binary content to file: %w", err)
	}

	return filePath, nil
}
