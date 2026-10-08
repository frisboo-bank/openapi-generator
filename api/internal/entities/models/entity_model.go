package models

import (
	"time"

	"github.com/google/uuid"
)

type Entity struct {
	EntityID    uuid.UUID  `db:"id"`
	Slug        string     `db:"slug"         query:"filter,order,search"`
	Name        string     `db:"name"         query:"order,search"`
	Description string     `db:"description"  query:"order,search"`
	VersionLock int64      `db:"version_lock"`
	DeletedAt    *time.Time `db:"deleted_at"    query:"filter"`
	CreatedAt   *time.Time `db:"created_at"   query:"filter,order"`
	UpdatedAt   *time.Time `db:"updated_at"   query:"filter,order"`
}
