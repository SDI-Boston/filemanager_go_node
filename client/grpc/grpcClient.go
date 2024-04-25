package client

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	pb "github.com/SDI-Boston/filemanager_go_node/client/proto"
	"google.golang.org/grpc"
)

func UploadClientFile() {
	serverAddr := "127.0.0.1:50051"
	filePath := "./grpc.txt"
	ownerID := "valdi"

	// Crear un contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	// Establecer la conexión utilizando DialContext
	conn, err := grpc.DialContext(
		ctx,
		serverAddr,
		grpc.WithInsecure(), // Usar conexión insegura para NGINX
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*1024)),
	)
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

	// Calcular el hash SHA256 del contenido del archivo
	hash := sha256.New()
	hash.Write(fileContent)
	//fileHash := fmt.Sprintf("%x", hash.Sum(nil))

	// Extraer el nombre del archivo con extensión
	fileName := filepath.Base(filePath)

	// Crear la solicitud de carga con el hash calculado
	uploadRequest := &pb.FileUploadRequest{
		FileId:     "99",
		OwnerId:    ownerID,
		BinaryFile: []byte(encodedContent),
		FileName:   fileName,
		//FileHash:   fileHash, // Incluir el hash calculado
	}

	// Abrir un flujo para enviar el archivo
	stream, err := client.Upload(ctx)
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

	// Imprimir la respuesta y la matriz de URLs
	fmt.Printf("File uploaded successfully. Path: %s\n", response)
}
