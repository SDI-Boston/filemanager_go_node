package server

import (
	"encoding/base64"
	"fmt"
	"os"

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

	// Subir archivo
	err = uploadToNFS(req)
	if err != nil {
		return fmt.Errorf("failed to upload file to NFS: %w", err)
	}

	// Respuesta
	filePath := fmt.Sprintf("172.171.240.20/files/%s/%s", req.OwnerId, req.FileId)
	err = stream.SendAndClose(&pb.FileUploadResponse{
		FileId: req.FileId,
		Urls:   []string{filePath},
	})
	if err != nil {
		return fmt.Errorf("failed to send upload response: %w", err)
	}

	return nil
}

func uploadToNFS(req *pb.FileUploadRequest) error {
	// Si el usuario nunca ha creado un archivo, crear un directorio para el usuario
	userPath := fmt.Sprintf("/mnt/nfs/%s", req.OwnerId)
	if _, err := os.Stat(userPath); os.IsNotExist(err) {
		err := os.Mkdir(userPath, 0755)
		if err != nil {
			return fmt.Errorf("failed to create user directory: %w", err)
		}
	}

	filePath := userPath + "/" + req.FileId
	fileUpload, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	defer fileUpload.Close()

	// Decodificar el contenido del binario que viene en base64
	decodedContent, err := base64.StdEncoding.DecodeString(string(req.BinaryFile))
	if err != nil {
		return fmt.Errorf("failed to decode binary content: %w", err)
	}

	// Escribir el contenido decodificado en el archivo
	_, err = fileUpload.Write(decodedContent)
	if err != nil {
		return fmt.Errorf("failed to write binary content to file: %w", err)
	}

	return nil
}
