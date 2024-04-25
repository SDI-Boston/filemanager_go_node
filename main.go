package main

import (
	"log"
	"net"
	"net/http"
	"time"

	pb "github.com/SDI-Boston/filemanager_go_node/proto"
	"github.com/SDI-Boston/filemanager_go_node/routes"
	"github.com/SDI-Boston/filemanager_go_node/server"
	"google.golang.org/grpc"
)

func main() {
	// Servidor gRPC
	grpcListener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer grpcListener.Close()

	// Servidor HTTP
	httpListener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer httpListener.Close()

	// Iniciar servidor gRPC con timeout de 5 minutos y limite de 1GB
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(1024*1024*1024),
		grpc.ConnectionTimeout(time.Minute*5),
	)
	pb.RegisterFileServiceServer(grpcServer, &server.FileService{})
	log.Println("gRPC server started")
	go func() {
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	router := routes.NewRouter()
	http.Handle("/", router)
	log.Println("HTTP server started")
	if err := http.Serve(httpListener, nil); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
