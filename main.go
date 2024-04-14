package main

import (
	"log"
	"net"

	"github.com/SDI-Boston/filemanager_go_node/client"
	pb "github.com/SDI-Boston/filemanager_go_node/proto"
	"github.com/SDI-Boston/filemanager_go_node/server"
	"google.golang.org/grpc"
)

func main() {
	//Levantar servidor gRPC
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//Crear servidor gRPC
	s := grpc.NewServer()

	//Registrar servicio
	pb.RegisterFileServiceServer(s, &server.FileService{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Println("Server started on port :50051")

	client.UploadClientFile()

}
