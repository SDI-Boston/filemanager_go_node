package server

import (
	"context"
	"fmt"
	"os"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) Upload(ctx context.Context, req *pb.FileUploadRequest) (*pb.FileUploadResponse, error) {
	// Save the binary file to NFS
	nfsPath := "/mnt/nfs"
	filePath := fmt.Sprintf("%s/%s", nfsPath, req.FileId)
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(req.BinaryFile)
	if err != nil {
		return nil, fmt.Errorf("failed to write file content: %w", err)
	}

	return &pb.FileUploadResponse{
		FileId:    req.FileId,
		FileRoute: filePath,
	}, nil
}
