package service

import (
	"context"
	"log"

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
	// filter := req.GetFilter()
	// log.Printf("receive a search-laptop request with filter: %v", filter)

	// err := server.laptopStore.Search(
	// 	stream.Context(),
	// 	filter,
	// 	func(laptop *pb.Laptop) error {
	// 		res := &pb.SearchLaptopResponse{Laptop: laptop}
	// 		err := stream.Send(res)
	// 		if err != nil {
	// 			return err
	// 		}

	// 		log.Printf("sent laptop with id: %s", laptop.GetId())
	// 		return nil
	// 	},
	// )

	// if err != nil {
	// 	return status.Errorf(codes.Internal, "unexpected error: %v", err)
	// }

	return nil
}
