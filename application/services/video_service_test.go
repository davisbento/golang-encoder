package services_test

import (
	"davisbento/golang-encoder/application/repositories"
	"davisbento/golang-encoder/application/services"
	"davisbento/golang-encoder/domain"
	"davisbento/golang-encoder/framework/database"
	dotenv "github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func init() {
	err := dotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}
}

func prepare() (*domain.Video, repositories.VideoRepositoryDb) {
	db := database.NewDBTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "convite.mp4"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{DbConn: db}

	return video, repo
}

func TestVideoServiceDownload(t *testing.T) {
	video, repo := prepare()

	service := services.NewVideoService()
	service.Video = video
	service.VideoRepository = repo

	err := service.Download("bucket-test-for-encoder")
	require.Nil(t, err)

	err = service.Fragment()
	require.Nil(t, err)

	err = service.Encode()
	require.Nil(t, err)

	err = service.Finish()
	require.Nil(t, err)
}
