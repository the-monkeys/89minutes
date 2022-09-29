package service

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/89minutes/89minutes/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		}
	} else {
		// Else assign a valid UUID
		newUUID, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new laptop id: %v", err)
		}

		story.Id = newUUID.String()
	}

	// Check for the context error
	if err := contextError(ctx); err != nil {
		return nil, err
	}

	story.CreateTime = timestamppb.Now()
	// Save the Story

	resp := &pb.CreateStoryResponse{
		Id: story.Id,
	}
	return resp, nil
}

func (server *StoryService) Update(context.Context, *pb.UpdateStoryRequest) (*pb.UpdateStoryResponse, error) {
	return nil, nil
}

func (server *StoryService) GetLatest(req *pb.GetLatestStoryRequest, stream pb.StoryService_GetLatestServer) error {
	log.Printf("receive a request to share latest reports")

	var wg sync.WaitGroup
	for i := 0; i < noOfLatestStories; i++ {
		wg.Add(1)
		go func(count int) {
			defer wg.Done()

			resp := pb.GetLatestStoryResponse{
				Story: &pb.Story{
					Id:    uuid.New().String(),
					Title: fmt.Sprintf("A New Article %d", count),
				},
			}

			if err := stream.Send(&resp); err != nil {
				log.Printf("send error %v", err)

			}
			log.Printf("finishing request number : %d", count)
		}(i)
	}

	wg.Wait()
	return nil
}
