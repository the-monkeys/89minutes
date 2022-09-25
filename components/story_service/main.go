package main

import (
	"log"

	"github.com/89minutes/89minutes/components/story_service/service"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("starting the gRPC service!")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("cannot load the environment variables: %s", err)
	}

	storyService := service.NewStoryService("OPENSEARCH")
	_ = storyService
}
