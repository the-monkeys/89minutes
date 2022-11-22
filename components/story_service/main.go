package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/89minutes/89minutes/components/story_service/osstore"
	"github.com/89minutes/89minutes/components/story_service/service"
	pb "github.com/89minutes/89minutes/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	opensearch "github.com/opensearch-project/opensearch-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	fileStorage = "cloud"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("cannot load the environment variables: %s", err)
	}

	osUser := os.Getenv("OSUSER")
	osPass := os.Getenv("OSPASS")
	osEndPoint := os.Getenv("OPENSEARCH")

	openClient := OpenClientClient(osEndPoint, osUser, osPass)
	osClient := osstore.NewOsClient(openClient)
	// newDiskStore := store.NewDiskFileStore("cloud")
	storyService := service.NewStoryService(osClient, *logrus.New(), fileStorage)
	grpcServer := grpc.NewServer()

	pb.RegisterStoryServiceServer(grpcServer, storyService)
	reflection.Register(grpcServer)

	address := fmt.Sprintf("0.0.0.0:%d", 8080)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("cannot create listener: %v", err)
	}

	log.Printf("gRPC server is starting at: %s", listener.Addr().String())
	go func() {
		err = grpcServer.Serve(listener)
		if err != nil {
			log.Fatalf("cannot start gRPC server: %v", err)
		}
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:8080",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	gwmux := runtime.NewServeMux()

	// Register Greeter
	err = pb.RegisterStoryServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}

func OpenClientClient(endPoint, user, pass string) *opensearch.Client {
	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: []string{endPoint},
		Username:  user,
		Password:  pass,
	})

	if err != nil {
		log.Fatalf("Cannot connect to opensearch: %v", err)
	}
	return client
}
