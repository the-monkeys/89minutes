package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/89minutes/89minutes/components/story_service/service"
	"github.com/89minutes/89minutes/components/story_service/store"
	pb "github.com/89minutes/89minutes/pb"
	"github.com/joho/godotenv"
	opensearch "github.com/opensearch-project/opensearch-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	newDiskStore := store.NewDiskFileStore("cloud")
	storyService := service.NewStoryService(openClient, newDiskStore)
	grpcServer := grpc.NewServer()

	pb.RegisterStoryServiceServer(grpcServer, storyService)
	reflection.Register(grpcServer)

	address := fmt.Sprintf("0.0.0.0:%d", 8080)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("cannot create listener: %v", err)
	}

	log.Printf("gRPC server is starting at: %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("cannot start gRPC server: %v", err)
	}

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
