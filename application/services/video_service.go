package services

import (
	storage "cloud.google.com/go/storage"
	"context"
	"davisbento/golang-encoder/application/repositories"
	"davisbento/golang-encoder/domain"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (service *VideoService) Download(bucketName string) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bucket := client.Bucket(bucketName)
	object := bucket.Object(service.Video.FilePath)

	reader, err := object.NewReader(ctx)
	if err != nil {
		return err
	}

	defer reader.Close()

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	path := os.Getenv("LOCAL_STORAGE_PATH") + "/" + service.Video.ID + ".mp4"
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = f.Write(body)
	if err != nil {
		return err
	}

	defer f.Close()
	log.Printf("Video %v downloaded with success!", service.Video.ID)
	return nil
}

func (service *VideoService) Fragment() error {
	path := os.Getenv("LOCAL_STORAGE_PATH") + "/" + service.Video.ID

	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		return err
	}

	source := os.Getenv("LOCAL_STORAGE_PATH") + "/" + service.Video.ID + ".mp4"
	target := path + "/" + service.Video.ID + ".frag"

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)
	log.Printf("Video %v fragmented with success!", service.Video.ID)
	return nil
}

func printOutput(output []byte) {
	if len(output) > 0 {
		log.Printf("=======> Output: %s\n", output)
	}
}
