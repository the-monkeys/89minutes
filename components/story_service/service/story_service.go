package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/89minutes/89minutes/pb"
	"github.com/opensearch-project/opensearch-go"
	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxImageSize = 1 << 20

type StoryService struct {
	OSClient *opensearch.Client
	logger   logrus.Logger
	storyDir string
	pb.UnimplementedStoryServiceServer
}

func NewStoryService(OSClient *opensearch.Client, logger logrus.Logger, storyDir string) *StoryService {
	return &StoryService{
		OSClient: OSClient,
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
	for {

		fileData, err = stream.Recv() //ignoring the data  TO-Do save files received

		if err != nil {
			if err == io.EOF {
				break
			}

			err = errors.Wrapf(err,
				"failed unexpectedly while reading chunks from stream")
			return
		}
		newPath := path.Join(server.storyDir, fileData.Id)
		// server.logger.Infof("The id: %+v", fileData)
		if err = os.MkdirAll(newPath, os.ModePerm); err != nil {
			server.logger.Errorf("cannot create the dir: %v", err)
			return err
		}

		if firstChunk { //first chunk contains file name

			if fileData.Filename != "" { //create file

				fp, err = os.Create(path.Join(newPath, filepath.Base(fileData.Filename)))

				if err != nil {
					server.logger.Errorf("Unable to create file %s, ERROR: %v", fileData.Filename, err)
					stream.SendAndClose(&pb.UploadStoryAndFilesRes{
						Message: "Unable to create file :" + fileData.Filename,
						Code:    pb.UploadStatusCode_Failed,
					})
					return
				}
				defer fp.Close()
			} else {
				server.logger.Errorf("FileName not provided in first chunk:  %s" + fileData.Filename)
				stream.SendAndClose(&pb.UploadStoryAndFilesRes{
					Message: "FileName not provided in first chunk:" + fileData.Filename,
					Code:    pb.UploadStatusCode_Failed,
				})
				return

			}
			filename = fileData.Filename
			firstChunk = false
		}

		err = writeToFp(fp, fileData.Content)
		if err != nil {
			server.logger.Errorf("Unable to write chunk of filename %s, ERROR: %v", fileData.Filename, err)
			stream.SendAndClose(&pb.UploadStoryAndFilesRes{
				Message: "Unable to write chunk of filename :" + fileData.Filename,
				Code:    pb.UploadStatusCode_Failed,
			})
			return
		}
	}

	//s.logger.Info().Msg("upload received")
	err = stream.SendAndClose(&pb.UploadStoryAndFilesRes{
		Message: "Upload received with success",
		Code:    pb.UploadStatusCode_Ok,
	})
	if err != nil {
		err = errors.Wrapf(err,
			"failed to send status code")
		return
	}
	fmt.Println("Successfully received and stored the file :" + filename + " in " + server.storyDir)

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
