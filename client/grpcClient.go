package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	pb "github.com/SDI-Boston/filemanager_go_node/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func UploadClientFile() {
	//serverAddr := "localhost:80"
	serverAddr := "localhost:50051"
	filePath := "./test/grpc.txt"
	ownerID := "1"

	// Crear las credenciales de transporte (en este caso, no seguras)
	creds := credentials.NewTLS(nil)

	// Establecer la conexión con las credenciales
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)

	// Leer el contenido del archivo
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Codificar el contenido del archivo a base64
	encodedContent := base64.StdEncoding.EncodeToString(fileContent)

	// Crear la solicitud de carga
	uploadRequest := &pb.FileUploadRequest{
		FileId:     "1",
		OwnerId:    ownerID,
		BinaryFile: []byte(encodedContent),
	}

	// Abrir un flujo para enviar el archivo
	stream, err := client.Upload(context.Background())
	if err != nil {
		log.Fatalf("Failed to open stream: %v", err)
	}

	// Enviar la solicitud de carga
	if err := stream.Send(uploadRequest); err != nil {
		log.Fatalf("Failed to send upload request: %v", err)
	}

	// Cerrar el flujo y recibir la respuesta
	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to receive response: %v", err)
	}

	fmt.Printf("File uploaded successfully. Path: %s\n", response.Urls[0])
}
