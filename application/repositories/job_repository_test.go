package repositories_test

import (
	"davisbento/golang-encoder/application/repositories"
	"davisbento/golang-encoder/domain"
	"davisbento/golang-encoder/framework/database"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJobRepositoryDbInsert(t *testing.T) {
	db := database.NewDBTest()
	defer db.Close()

	video, err := insertVideoTest()

	require.Nil(t, err)
	require.NotEmpty(t, video.ID)

	job, err := domain.NewJob("Transcoding", "Status", video)
	require.Nil(t, err)

	repo := repositories.NewJobRepository(db)
	repo.Insert(job)

	j, err := repo.Find(job.ID)
	require.NotEmpty(t, j.ID)
	require.Nil(t, err)
}

func TestJobDbUpdate(t *testing.T) {
	db := database.NewDBTest()
	defer db.Close()

	video, err := insertVideoTest()

	require.Nil(t, err)
	require.NotEmpty(t, video.ID)

	job, err := domain.NewJob("Transcoding", "Pending", video)
	require.Nil(t, err)

	repo := repositories.NewJobRepository(db)
	repo.Insert(job)

	job.Status = "Completed"
	job, err = repo.Update(job)

	require.Nil(t, err)
	require.Equal(t, "Completed", job.Status)
}
