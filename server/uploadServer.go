package server

import (
	"context"
	"fmt"
	"os"

	pb "github.com/SDI-Boston/filemanager_go_node/proto"
)

type FileService struct{}

func (s *FileService) Upload(ctx context.Context, req *pb.FileUploadRequest) (*pb.FileUploadResponse, error) {

	// Subir archivo
	err := uploadToNFS(req)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to NFS: %w", err)
	}

	// Respuesta
	filePath := fmt.Sprintf("/mnt/nfs/%s/%s", req.OwnerId, req.FileId)
	return &pb.FileUploadResponse{
		FileId: req.FileId,
		Urls:   []string{filePath},
	}, nil
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
		return fmt.Errorf("failed to create test file: %w", err)
	}
	defer fileUpload.Close()

	// Escribir el contenido del archivo binario en el archivo
	_, err = fileUpload.Write(req.BinaryFile)
	if err != nil {
		return fmt.Errorf("failed to write binary content to file: %w", err)
	}

	return nil
}
