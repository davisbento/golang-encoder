package repositories_test

import (
	"davisbento/golang-encoder/application/repositories"
	"davisbento/golang-encoder/domain"
	"davisbento/golang-encoder/framework/database"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func insertVideoTest() (*domain.Video, error) {
	db := database.NewDBTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path/to/file.mp4"
	video.CreatedAt = time.Now()
	repo := repositories.NewVideoRepository(db)
	repo.Insert(video)

	v, err := repo.Find(video.ID)

	if err != nil {
		return nil, err
	}

	return v, nil
}

func TestVideoRepositoryDBInsert(t *testing.T) {
	v, err := insertVideoTest()

	require.NotEmpty(t, v.ID)
	require.Nil(t, err)
	require.Equal(t, "path/to/file.mp4", v.FilePath)
}
