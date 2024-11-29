package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/gerbsec/Helheim/proto" // Import the generated protobuf package
)

func main() {
	// Load the CA certificate
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("cert/server.crt")
	if err != nil {
		log.Fatalf("could not read server certificate: %v", err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("failed to append server certificate")
	}

	// Create TLS credentials using the server's certificate
	creds := credentials.NewTLS(&tls.Config{
		RootCAs: certPool,
	})

	// Connect to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewCommandControlClient(conn)

	fmt.Println("Connected to the C2 server. Type 'status' to get server status, 'exec <command>' to execute a command, or 'exit' to quit.")
	reader := bufio.NewReader(os.Stdin)

	// Start interactive session
	for {
		fmt.Print("> ") // Command prompt

		// Read user input
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("could not read input: %v", err)
		}

		// Trim whitespace and newline characters
		input = strings.TrimSpace(input)

		// Handle the 'exit' command
		if input == "exit" {
			fmt.Println("Exiting session...")
			break
		}

		// Split input into command and arguments
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue // No command entered
		}

		command := parts[0]

		// Switch based on the command
		switch command {
		case "status":
			getStatus(client)
		case "exec":
			if len(parts) < 2 {
				fmt.Println("usage: exec <command>")
				continue
			}
			execCommand(client, strings.Join(parts[1:], " "))
		default:
			fmt.Printf("Unknown command: %s\n", command)
		}
	}
}

// Get the status from the server
func getStatus(client pb.CommandControlClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.StatusRequest{}
	res, err := client.GetStatus(ctx, req)
	if err != nil {
		log.Printf("could not get status: %v", err)
		return
	}

	fmt.Printf("Server status: %s\n", res.Status)
}

// Execute a command on the server
func execCommand(client pb.CommandControlClient, command string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.CommandRequest{
		Command: command,
	}

	res, err := client.ExecuteCommand(ctx, req)
	if err != nil {
		log.Printf("could not execute command: %v", err)
		return
	}

	fmt.Printf("Command result: %s\n", res.Result)
}
