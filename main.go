package main

import (
	"log"
	"net"
	"net/http"

	pb "github.com/SDI-Boston/filemanager_go_node/proto"
	"github.com/SDI-Boston/filemanager_go_node/routes"
	"github.com/SDI-Boston/filemanager_go_node/server"
	"google.golang.org/grpc"
)

func main() {
	// Levantar servidor gRPC
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterFileServiceServer(s, &server.FileService{})

	// Iniciar servidor gRPC en una goroutine
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	log.Println("Server started on port :50051")

	// Iniciar servidor HTTP para descargas
	router := routes.NewRouter()
	http.Handle("/", router)
	log.Println("HTTP download server started on port :8080")
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	select {}

}
