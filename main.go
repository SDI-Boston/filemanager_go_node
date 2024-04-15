package main

import (
	"log"
	"net"
	"net/http"
	"strconv"

	pb "github.com/SDI-Boston/filemanager_go_node/proto"
	"github.com/SDI-Boston/filemanager_go_node/routes"
	"github.com/SDI-Boston/filemanager_go_node/server"
	"google.golang.org/grpc"
)

func main() {
	// Rango de puertos para escalabilidad
	for port := 50051; port <= 50060; port++ {
		lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
		if err == nil {
			// Found an available port, start serving
			log.Printf("Server listening on port %d\n", port)
			startServers(lis)
			return
		}
	}

	log.Fatalf("Failed to find available port in range")
}

func startServers(lis net.Listener) {
	s := grpc.NewServer()
	pb.RegisterFileServiceServer(s, &server.FileService{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	log.Println("gRPC server started")

	router := routes.NewRouter()
	http.Handle("/", router)
	log.Println("HTTP server started")
	if err := http.Serve(lis, nil); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
