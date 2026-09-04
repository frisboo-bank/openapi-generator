package dtos

import "time"

type EntityDto struct {
	EntityID    string
	Slug        string
	Name        string
	Description *string
	VersionLock int32
	DeletedAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
