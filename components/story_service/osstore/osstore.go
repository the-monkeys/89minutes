package osstore

import (
	"context"
	"fmt"
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

func (client *OsClient) CreateNewStory(ip *pb.UploadStoryAndFilesReq, filePath string, index string) error {
	story := StoryToString(ip, filePath)
	// logrus.Info("Story: ", story)

	document := strings.NewReader(story)

	req := opensearchapi.IndexRequest{
		Index:      index,
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

func StoryToString(ip *pb.UploadStoryAndFilesReq, filePath string) string {
	return fmt.Sprintf(`{
		"id":         "%s",
		"title":      "%s",
		"author":     "%s",
		"filepath":   "%s",
		"isDraft":    "%v",
		"created_at": "%v"
	}`, ip.Id, ip.Title, ip.Author, filePath, ip.IsDraft, time.Now())
}
