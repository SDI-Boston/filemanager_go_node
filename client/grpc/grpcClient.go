package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	pb "github.com/SDI-Boston/filemanager_go_node/client/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func UploadClientFile() {
	serverAddr := "node.eastus.cloudapp.azure.com:5000"
	filePath := "./grpc.txt"
	ownerID := "owner1"

	// Establishing an insecure connection
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)

	// Read the content of the file
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Encode the file content to base64
	encodedContent := base64.StdEncoding.EncodeToString(fileContent)

	// Create the upload request
	uploadRequest := &pb.FileUploadRequest{
		FileId:     "file1",
		OwnerId:    ownerID,
		BinaryFile: []byte(encodedContent),
	}

	// Open a stream to send the file
	stream, err := client.Upload(context.Background())
	if err != nil {
		log.Fatalf("Failed to open stream: %v", err)
	}

	// Send the upload request
	if err := stream.Send(uploadRequest); err != nil {
		log.Fatalf("Failed to send upload request: %v", err)
	}

	// Close the stream and receive the response
	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to receive response: %v", err)
	}

	fmt.Printf("File uploaded successfully. Path: %s\n", response.Urls[0])
}
