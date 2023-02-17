package domain

import "time"

type Job struct {
	ID               string
	OutputBucketPath string
	Status           string
	Video            *Video
	CreatedAt        time.Time
	UpdateAt         time.Time
}
