package store

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
)

type FileStore interface {
	Save(storyId string, fileType string, fileData bytes.Buffer) (string, error)
}

type DiskFileStore struct {
	mutex   sync.RWMutex
	fileDir string
	file    map[string]*FileInfo
}

type FileInfo struct {
	StoryId string
	Type    string
	Path    string
}

func NewDiskFileStore(fileDir string) *DiskFileStore {
	return &DiskFileStore{
		fileDir: fileDir,
		file:    make(map[string]*FileInfo),
	}
}

func (store *DiskFileStore) Save(storyId string, fileType string, fileData bytes.Buffer) (string, error) {
	imageID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("cannot generate file id: %w", err)
	}
	log.Println("****storyId: ", storyId)
	log.Println("****fileType: ", fileType)
	// log.Println("****FILE PATH: ", imagePath)

	imagePath := fmt.Sprintf("%s/%s", store.fileDir, imageID.String()+fileType)

	log.Println("****FILE PATH: ", imagePath)

	file, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("cannot create file file: %w", err)
	}

	_, err = fileData.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("cannot write file to file: %w", err)
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.file[imageID.String()] = &FileInfo{
		StoryId: storyId,
		Type:    fileType,
		Path:    imagePath,
	}

	return imageID.String(), nil
}
