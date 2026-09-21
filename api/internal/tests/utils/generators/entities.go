package generators

import (
	"time"

	"frisboo-bank/openapi-generator-service/internal/entities/models"

	"pgregory.net/rapid"
)

func GenerateEntityToSave() *rapid.Generator[*models.Entity] {
	return rapid.Custom(func(t *rapid.T) *models.Entity {
		description := GenerateDescription().Draw(t, "description")

		return &models.Entity{
			Slug:        GenerateSlug().Draw(t, "slug"),
			Name:        GenerateName().Draw(t, "name"),
			Description: description,
		}
	})
}

func GenerateEntity() *rapid.Generator[*models.Entity] {
	return rapid.Custom(func(t *rapid.T) *models.Entity {
		createdAt := GenerateTimestamp(10, false).Draw(t, "createdAt")
		updatedAt := GenerateTimestamp(10, false).Draw(t, "updatedAt")
		var hiddenAt *time.Time
		if rapid.Bool().Draw(t, "hasHiddenAt") {
			hiddenAt = updatedAt
		}

		e := GenerateEntityToSave().Draw(t, "entity")
		e.EntityID = GenerateEntityID().Draw(t, "entityID")
		e.VersionLock = GenerateVersion().Draw(t, "versionLock")
		e.HiddenAt = hiddenAt
		e.CreatedAt = createdAt
		e.UpdatedAt = updatedAt

		return e
	})
}
