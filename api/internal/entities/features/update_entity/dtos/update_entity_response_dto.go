package dtos

import "time"

type UpdateEntityResponseDto struct {
	Slug        string
	Name        string
	Description string
	VersionLock int64
	HiddenAt    *time.Time
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
