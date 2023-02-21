package services

import (
	"cloud.google.com/go/storage"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type UploadManager struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewUploadManager() UploadManager {
	return UploadManager{}
}

func (manager *UploadManager) UploadObject(objectPath string, client *storage.Client, ctx context.Context) error {
	// path/x/bucket/video.mp4
	// split: [path/x/bucket, video.mp4]
	path := strings.Split(objectPath, os.Getenv("LOCAL_STORAGE_PATH")+"/")

	f, err := os.Open(objectPath)
	if err != nil {
		return err
	}

	defer f.Close()

	wc := client.Bucket(manager.OutputBucket).Object(path[1]).NewWriter(ctx)
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err = io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

func (manager *UploadManager) LoadPaths() error {
	err := filepath.Walk(manager.VideoPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			manager.Paths = append(manager.Paths, path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func getClientUpload() (*storage.Client, context.Context, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}
