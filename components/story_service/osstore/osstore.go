package osstore

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/89minutes/89minutes/pb"
	opensearch "github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"github.com/sirupsen/logrus"
)

type Story struct {
	id         string
	title      string
	author     string
	filepath   string
	isDraft    bool
	created_at time.Time
	updated_at time.Time
}

type OsClient struct {
	client *opensearch.Client
}

func NewOsClient(client *opensearch.Client) *OsClient {
	return &OsClient{
		client: client,
	}
}

func (client *OsClient) CreateNewStory(ip *pb.UploadStoryAndFilesReq, filePath string) error {
	story := NewStory(ip, filePath)
	byteSlice, err := json.Marshal(story)
	if err != nil {
		return err
	}

	logrus.Infof("The byteslice: %v", string(byteSlice))
	document := strings.NewReader(string(byteSlice))

	req := opensearchapi.IndexRequest{
		Index:      "go-test-index1",
		DocumentID: ip.Id,
		Body:       document,
	}
	insertResponse, err := req.Do(context.Background(), client.client)
	if err != nil {
		logrus.Errorf("cannot create a new document, error: %v", err)
		return err
	}

	logrus.Infof("insert response: %+v", insertResponse)

	return nil
}

func NewStory(ip *pb.UploadStoryAndFilesReq, filePath string) *Story {
	return &Story{
		id:         ip.Id,
		title:      ip.Title,
		author:     ip.Author,
		filepath:   filePath,
		isDraft:    ip.IsDraft,
		created_at: time.Now(),
	}
}
