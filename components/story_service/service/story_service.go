package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/89minutes/89minutes/components/story_service/store"
	"github.com/89minutes/89minutes/pb"
	"github.com/google/uuid"
	"github.com/opensearch-project/opensearch-go"
	opensearchapi "github.com/opensearch-project/opensearch-go/opensearchapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type StoryService struct {
	OSClient  *opensearch.Client
	fileStore store.FileStore
	pb.UnimplementedStoryServiceServer
}

func NewStoryService(OSClient *opensearch.Client, fileStore store.FileStore) *StoryService {
	return &StoryService{
		OSClient:  OSClient,
		fileStore: fileStore,
	}
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

	log.Printf("The files: %v", story.FileContent)
	// TODO: store files
	fileData := bytes.Buffer{}

	for _, file := range story.FileContent {
		_, err := fileData.Write(file)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot write chunk data: %v", err)
		}

		fileId, err := server.fileStore.Save(story.Id, "image", fileData)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot save image to the store: %v", err)
		}

		// Replace story filenplaceholder with storyidn.ext
		fmt.Printf("File Id: %s", fileId)
	}

	// TODO: Store special text

	// Add or update tags
	tags := story.Tag
	log.Printf("Tags are: %v", tags)

	story.CreateTime = timestamppb.Now()

	storyBytesSlice, err := json.Marshal(story)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot marshal the story: %v", err)
	}

	newReader := strings.NewReader(string(storyBytesSlice))

	// Insert into opensearch
	OSReq := opensearchapi.IndexRequest{
		Index: storyIndex,
		Body:  newReader,
	}

	insertResponse, err := OSReq.Do(context.Background(), server.OSClient)
	if err != nil {
		fmt.Println("failed to insert document ", err)
		os.Exit(1)
	}
	fmt.Println("Inserting a document")
	fmt.Println(insertResponse)

	resp := &pb.CreateStoryResponse{
		Id: story.Id,
	}
	return resp, nil
}

func (server *StoryService) Update(ctx context.Context, req *pb.UpdateStoryRequest) (*pb.UpdateStoryResponse, error) {
	id := req.GetId()
	title := req.GetTitle()
	textContent := req.GetTextContent()
	updateTime := timestamppb.Now()

	log.Printf("Updating story %s, titled: %s", id, title)

	// TODO: Check if the story exits

	// TODO: Update the story
	_, _ = textContent, updateTime

	resp := &pb.UpdateStoryResponse{
		Id: id,
	}
	return resp, nil
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
