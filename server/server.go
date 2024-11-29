package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/gerbsec/Helheim/proto" // Update this import to match your module path
)

// Server represents the gRPC server
type Server struct {
	pb.UnimplementedCommandControlServer
}

// ExecuteCommand handles the execution of a command
func (s *Server) ExecuteCommand(ctx context.Context, req *pb.CommandRequest) (*pb.CommandResponse, error) {
	log.Printf("Received command: %s\n", req.Command)
	// Here we would execute the command, but for now, we'll return a dummy response
	result := fmt.Sprintf("Executed command: %s", req.Command)
	return &pb.CommandResponse{Result: result}, nil
}

// GetStatus returns the current status of the server
func (s *Server) GetStatus(ctx context.Context, req *pb.StatusRequest) (*pb.StatusResponse, error) {
	status := "Server is up and running"
	return &pb.StatusResponse{Status: status}, nil
}

func main() {
	// Load server certificates from disk
	cert, err := tls.LoadX509KeyPair("cert/server.crt", "cert/server.key")
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}

	// Create TLS credentials
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
	})

	// Set up a TCP listener on port 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server instance with TLS credentials
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterCommandControlServer(grpcServer, &Server{})

	log.Println("C2 server is running on port 50051...")

	// Start the server
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
