package service

import (
	"context"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/89minutes/89minutes/components/story_service/osstore"
	"github.com/89minutes/89minutes/pb"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxImageSize = 1 << 20

type StoryService struct {
	// OSClient *opensearch.Client
	OsClient *osstore.OsClient
	logger   logrus.Logger
	storyDir string
	pb.UnimplementedStoryServiceServer
}

func NewStoryService(OsClient *osstore.OsClient, logger logrus.Logger, storyDir string) *StoryService {
	return &StoryService{
		OsClient: OsClient,
		logger:   logger,
		storyDir: storyDir,
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

func (server *StoryService) UploadStoryAndFiles(stream pb.StoryService_UploadStoryAndFilesServer) (err error) {
	firstChunk := true

	var fp *os.File
	var fileData *pb.UploadStoryAndFilesReq
	var filename string
	// var flag = true

	var info *pb.UploadStoryAndFilesReq

	for {
		fileData, err = stream.Recv() //ignoring the data  TO-Do save files received
		if err == io.EOF {
			break
		}
		if err != nil {
			server.logger.Errorf("failed unexpectedly while reading chunks, error: %v", err)
			return err
		}

		newPath := path.Join(server.storyDir, fileData.Id)

		if err = os.MkdirAll(newPath, os.ModePerm); err != nil {
			server.logger.Errorf("cannot create the dir: %v", err)
			return err
		}

		info = fileData
		// Store information into opensearch one time if flag it true
		// if flag {
		// 	logrus.Info("The flag is: ", flag)
		// 	if err := server.OsClient.CreateNewStory(fileData, newPath); err != nil {
		// 		return err
		// 	}
		// 	flag = false
		// }

		info = fileData

		if firstChunk { //first chunk contains file name
			if fileData.Filename != "" { //create file

				fp, err = os.Create(path.Join(newPath, filepath.Base(fileData.Filename)))
				if err != nil {
					server.logger.Errorf("Unable to create file %s, ERROR: %v", fileData.Filename, err)
					stream.SendAndClose(&pb.UploadStoryAndFilesRes{Message: "Unable to create file :" + fileData.Filename, Code: pb.UploadStatusCode_Failed})
					return
				}

				defer fp.Close()
			} else {
				server.logger.Errorf("FileName not provided in first chunk:  %s" + fileData.Filename)
				stream.SendAndClose(&pb.UploadStoryAndFilesRes{Message: "FileName not provided in first chunk:" + fileData.Filename, Code: pb.UploadStatusCode_Failed})
				return
			}

			filename = fileData.Filename
			firstChunk = false
		}

		err = writeToFp(fp, fileData.Content)
		if err != nil {
			server.logger.Errorf("Unable to write chunk of filename %s, ERROR: %v", fileData.Filename, err)
			stream.SendAndClose(&pb.UploadStoryAndFilesRes{Message: "Unable to write chunk of filename :" + fileData.Filename, Code: pb.UploadStatusCode_Failed})
			return
		}
	}

	err = stream.SendAndClose(&pb.UploadStoryAndFilesRes{Message: "Upload received with success", Code: pb.UploadStatusCode_Ok})
	if err != nil {
		server.logger.Errorf("failed to send status code, error: %v", err)
		return
	}

	server.logger.Infof("Successfully received and stored the file: %s, in  dir: %s", filename, server.storyDir)

	//
	//  This needs to be changed
	//
	if err := server.OsClient.CreateNewStory(info, "newPath"); err != nil {
		return err
	}
	return nil
}

//writeToFp takes in a file pointer and byte array and writes the byte array into the file
//returns error if pointer is nil or error in writing to file
func writeToFp(fp *os.File, data []byte) error {
	w := 0
	n := len(data)
	for {

		nw, err := fp.Write(data[w:])
		if err != nil {
			return err
		}
		w += nw
		if nw >= n {
			return nil
		}
	}

}
