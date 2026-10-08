package dtos

import "time"

type CreateEntityResponseDto struct {
	Slug        string
	Name        string
	Description *string
	VersionLock int64
	DeletedAt    *time.Time
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
