package service

import (
	"context"
	"log"

	"github.com/89minutes/89minutes/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StoryService struct {
	OSClient string
}

func NewStoryService(OSClient string) *StoryService {
	return &StoryService{OSClient}
}

func (server *StoryService) Create(ctx context.Context, req *pb.CreateStoryRequest) (*pb.CreateStoryResponse, error) {
	story := req.GetStory()
	log.Printf("received a create-story request with id: %s", story.Id)

	// Check if Id exists
	if len(story.Id) > 0 {
		// If Id exists then is it a valid UUID
		_, err := uuid.Parse(story.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "story id is not a valid UUID: %v", err)
		} else {
			// Else assign a valid UUID
			id, err := uuid.NewRandom()
			if err != nil {
				return nil, status.Errorf(codes.Internal, "cannot generate a new laptop id: %v", err)
			}

			story.Id = id.String()
		}
	}

	// Check for the context error
	if err := contextError(ctx); err != nil {
		return nil, err
	}

	// Save the Story

	resp := &pb.CreateStoryResponse{
		Id: story.Id,
	}
	return resp, nil
}
