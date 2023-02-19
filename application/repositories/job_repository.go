package repositories

import (
	"davisbento/golang-encoder/domain"
	"github.com/jinzhu/gorm"
)

type JobRepository interface {
	Insert(job *domain.Job) (*domain.Job, error)
	Find(id string) (*domain.Job, error)
	Update(job *domain.Job) (*domain.Job, error)
}

type JobRepositoryDb struct {
	DbConn *gorm.DB
}

func NewJobRepository(db *gorm.DB) *JobRepositoryDb {
	return &JobRepositoryDb{DbConn: db}
}

func (repo JobRepositoryDb) Insert(job *domain.Job) (*domain.Job, error) {
	err := repo.DbConn.Create(job).Error

	if err != nil {
		return nil, err
	}

	return job, nil
}

func (repo JobRepositoryDb) Find(id string) (*domain.Job, error) {
	var job domain.Job

	repo.DbConn.Preload("Video").First(&job, "id = ?", id)

	if job.ID == "" {
		return nil, nil
	}

	return &job, nil
}

func (repo JobRepositoryDb) Update(job *domain.Job) (*domain.Job, error) {
	err := repo.DbConn.Save(&job).Error

	if err != nil {
		return nil, err
	}

	return job, nil
}
