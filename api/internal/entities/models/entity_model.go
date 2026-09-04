package models

import "time"

type Entity struct {
	EntityID    string     `db:"id"           json:"id"`
	Slug        string     `db:"slug"         json:"slug"`
	Name        string     `db:"name"         json:"name"`
	Description *string    `db:"description"  json:"description,omitempty"`
	VersionLock int32      `db:"version_lock" json:"version_lock"`
	HiddenAt    *time.Time `db:"hidden_at"    json:"hidden_at"`
	CreatedAt   time.Time  `db:"created_at"   json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"   json:"updated_at"`
}
