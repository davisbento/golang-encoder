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
	err := os.Mkdir(os.Getenv("LOCAL_STORAGE_PATH")+"/"+service.Video.ID, os.ModePerm)
	if err != nil {
		return err
	}

	source := os.Getenv("LOCAL_STORAGE_PATH") + "/" + service.Video.ID + ".mp4"
	target := os.Getenv("LOCAL_STORAGE_PATH") + "/" + service.Video.ID + ".frag"

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (service *VideoService) Encode() error {
	var cmdArgs []string
	cmdArgs = append(cmdArgs, os.Getenv("LOCAL_STORAGE_PATH")+"/"+service.Video.ID+".frag")
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, os.Getenv("LOCAL_STORAGE_PATH")+"/"+service.Video.ID)
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin/")
	cmd := exec.Command("mp4dash", cmdArgs...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	printOutput(output)

	return nil
}

func (service *VideoService) Finish() error {

	err := os.Remove(os.Getenv("LOCAL_STORAGE_PATH") + "/" + service.Video.ID + ".mp4")
	if err != nil {
		log.Println("error removing mp4 ", service.Video.ID, ".mp4")
		return err
	}

	err = os.Remove(os.Getenv("LOCAL_STORAGE_PATH") + "/" + service.Video.ID + ".frag")
	if err != nil {
		log.Println("error removing frag ", service.Video.ID, ".frag")
		return err
	}

	err = os.RemoveAll(os.Getenv("LOCAL_STORAGE_PATH") + "/" + service.Video.ID)
	if err != nil {
		log.Println("error removing mp4 ", service.Video.ID, ".mp4")
		return err
	}

	log.Println("files have been removed: ", service.Video.ID)

	return nil

}

func printOutput(output []byte) {
	if len(output) > 0 {
		log.Printf("=======> Output: %s\n", output)
	}
}
