package main

import (
	"fmt"
	"log"
	"net"

	"github.com/89minutes/89minutes/components/story_service/service"
	"github.com/89minutes/89minutes/pb"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("cannot load the environment variables: %s", err)
	}

	storyService := service.NewStoryService("OPENSEARCH")
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
