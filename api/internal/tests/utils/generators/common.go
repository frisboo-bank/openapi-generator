package generators

import (
	"time"

	"github.com/google/uuid"
	"pgregory.net/rapid"
)

func GenerateEntityID() *rapid.Generator[uuid.UUID] {
	return rapid.Custom(func(t *rapid.T) uuid.UUID {
		return uuid.New()
	})
}

func GenerateSlug() *rapid.Generator[string] {
	return rapid.StringMatching(`^[a-z][a-z0-9_-]{2,49}$`)
}

func GenerateName() *rapid.Generator[string] {
	return rapid.StringN(3, 100, -1)
}

func GenerateDescription() *rapid.Generator[string] {
	return rapid.StringN(0, 500, -1)
}

func GenerateVersion() *rapid.Generator[int64] {
	return rapid.Int64Range(1, 100)
}

type Timestamps struct {
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func GenerateTimestamps(maxCreatedYears, maxUpdateDeltaDays int, updatedAtNullable bool) *rapid.Generator[Timestamps] {
	return rapid.Custom(func(t *rapid.T) Timestamps {
		createdAt := GenerateTimestamp(maxCreatedYears, false).Draw(t, "createdAt")

		if updatedAtNullable && !rapid.Bool().Draw(t, "hasUpdatedAt") {
			return Timestamps{
				CreatedAt: createdAt,
				UpdatedAt: nil,
			}
		}

		maxUpdateDeltaSeconds := int64(maxUpdateDeltaDays) * 24 * 3600
		updatedAtDeltaSeconds := rapid.Int64Range(0, maxUpdateDeltaSeconds).Draw(t, "updatedAtDeltaSeconds")
		updatedAt := createdAt.Add(time.Duration(updatedAtDeltaSeconds) * time.Second)

		return Timestamps{
			CreatedAt: createdAt,
			UpdatedAt: &updatedAt,
		}
	})
}

func GenerateTimestamp(maxYears int, nullable bool) *rapid.Generator[*time.Time] {
	return rapid.Custom(func(t *rapid.T) *time.Time {
		if nullable && rapid.Bool().Draw(t, "isNull") {
			return nil
		}

		maxAge := time.Duration(maxYears) * 365 * 24 * time.Hour
		ago := rapid.Int64Range(0, int64(maxAge.Seconds())).Draw(t, "secondsAgo")
		tm := time.Now().UTC().Add(-time.Duration(ago) * time.Second)

		return &tm
	})
}
