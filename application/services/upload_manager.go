package services

import (
	"cloud.google.com/go/storage"
	"context"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"runtime"
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

func (manager *UploadManager) ProcessUpload(concurrency int, doneUpload chan string) error {
	in := make(chan int, runtime.NumCPU()) // channel to send the index of the paths
	returnChannel := make(chan string)     // channel to receive the result of the upload

	err := manager.LoadPaths()
	if err != nil {
		return err
	}

	client, ctx, err := getClientUpload()
	if err != nil {
		return err
	}

	// start the workers with goroutines
	for process := 0; process < concurrency; process++ {
		go manager.uploadWorker(in, returnChannel, client, ctx)
	}

	// send the index of the paths to the workers
	go func() {
		for x := 0; x < len(manager.Paths); x++ {
			in <- x
		}
		close(in)
	}()

	// keep listening the return channel to check if any error occurs
	for r := range returnChannel {
		if r != "success" {
			// if any error occurs, we stop the upload
			doneUpload <- r
			log.Println("Upload failed")
			break
		}

		log.Println("Upload success")
	}

	return nil
}

func (manager *UploadManager) uploadWorker(in chan int, returnChannel chan string, uploadClient *storage.Client, ctx context.Context) {
	for x := range in {
		err := manager.UploadObject(manager.Paths[x], uploadClient, ctx)
		if err != nil {
			manager.Errors = append(manager.Errors, manager.Paths[x])
			log.Printf("Error uploading file %v. Error: %v", manager.Paths[x], err)
			returnChannel <- err.Error()
		}

		returnChannel <- "success"
	}

	returnChannel <- "finished"
}

func getClientUpload() (*storage.Client, context.Context, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}
