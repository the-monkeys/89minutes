package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"

	"github.com/89minutes/89minutes/components/story_service/store"
	"github.com/89minutes/89minutes/pb"
	"github.com/google/uuid"
	"github.com/opensearch-project/opensearch-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const maxImageSize = 1 << 20

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

func (server *StoryService) PingPong(ctx context.Context, req *pb.Ping) (*pb.Pong, error) {
	if req.Ping != "ping" {
		return nil, status.Errorf(codes.InvalidArgument, "The request is not ping")
	}
	return &pb.Pong{
		Ping: "Pong",
	}, nil
}

func (server *StoryService) GetBlob(stream pb.StoryService_GetBlobServer) error {
	req, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot receive image info")
	}

	fileSize := req.GetInfo().GetFileSize()
	filePath := req.GetInfo().GetName()
	fileId := uuid.New().String()
	_, fileName := filepath.Split(filePath)

	log.Printf("receive an request for chunk %d filename %s", fileSize, fileName)

	fileData := bytes.Buffer{}
	fileSizeNew := 0

	for {
		err := contextError(stream.Context())
		if err != nil {
			return err
		}

		log.Print("waiting to receive more data")

		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot receive chunk data: %v", err)
		}

		chunk := req.GetChunkData()
		size := len(chunk)

		log.Printf("received a chunk with size: %d", size)

		fileSizeNew += size
		if fileSize > maxImageSize {
			return status.Errorf(codes.InvalidArgument, "image is too large: %d > %d", fileSize, maxImageSize)
		}

		// write slowly
		// time.Sleep(time.Second)

		_, err = fileData.Write(chunk)
		if err != nil {
			return status.Errorf(codes.Internal, "cannot write chunk data: %v", err)
		}
	}

	fileID, err := server.fileStore.Save(fileId, fileName, fileData)
	if err != nil {
		return status.Errorf(codes.Internal, "cannot save image to the store: %v", err)
	}

	res := &pb.StoryResponse{
		Id:   fileID,
		Size: uint32(fileSize),
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot send response: %v", err)
	}

	log.Printf("saved image with id: %s, size: %d", fileID, fileSize)
	return nil
}

func (server *StoryService) Create(ctx context.Context, req *pb.CreateStoryRequest) (*pb.CreateStoryResponse, error) {
	// story := req.GetStory()
	byteX, err := json.MarshalIndent(req, "", "    ")
	if err != nil {
		fmt.Printf("Error while marshalling json: %v", err)
		// return nil, err
	}

	ioutil.WriteFile("test/sample_io.json", byteX, 0666)

	log.Println("Hello")
	// story := req.GetStory()
	// log.Printf("received a create-story request with id: %s", story.Id)

	// // Check if Id exists
	// if len(story.Id) > 0 {
	// 	// If Id exists then is it a valid UUID
	// 	_, err := uuid.Parse(story.Id)
	// 	if err != nil {
	// 		return nil, status.Errorf(codes.InvalidArgument, "story id is not a valid UUID: %v", err)
	// 	}
	// } else {
	// 	// Else assign a valid UUID
	// 	newUUID, err := uuid.NewRandom()
	// 	if err != nil {
	// 		return nil, status.Errorf(codes.Internal, "cannot generate a new laptop id: %v", err)
	// 	}

	// 	story.Id = newUUID.String()
	// }

	// // Check for the context error
	// if err := contextError(ctx); err != nil {
	// 	return nil, err
	// }

	// log.Printf("The files: %v", story.FileContent)
	// // TODO: store files
	// fileData := bytes.Buffer{}

	// for _, file := range story.FileContent {
	// 	_, err := fileData.Write(file)
	// 	if err != nil {
	// 		return nil, status.Errorf(codes.Internal, "cannot write chunk data: %v", err)
	// 	}

	// 	fileId, err := server.fileStore.Save(story.Id, "image", fileData)
	// 	if err != nil {
	// 		return nil, status.Errorf(codes.Internal, "cannot save image to the store: %v", err)
	// 	}

	// 	// Replace story filenplaceholder with storyidn.ext
	// 	fmt.Printf("File Id: %s", fileId)
	// }

	// // TODO: Store special text

	// // Add or update tags
	// tags := story.Tag
	// log.Printf("Tags are: %v", tags)

	// story.CreateTime = timestamppb.Now()

	// storyBytesSlice, err := json.Marshal(story)
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "cannot marshal the story: %v", err)
	// }

	// newReader := strings.NewReader(string(storyBytesSlice))

	// // Insert into opensearch
	// OSReq := opensearchapi.IndexRequest{
	// 	Index: storyIndex,
	// 	Body:  newReader,
	// }

	// insertResponse, err := OSReq.Do(context.Background(), server.OSClient)
	// if err != nil {
	// 	fmt.Println("failed to insert document ", err)
	// 	os.Exit(1)
	// }
	// fmt.Println("Inserting a document")
	// fmt.Println(insertResponse)

	// resp := &pb.CreateStoryResponse{
	// 	Id: story.Id,
	// }
	return nil, nil
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
