package models

type AppOptions struct {
	Name string `mapstructure:"name" json:"name"`
	Version string `mapstructure:"version" json:"version"`
	Description string `mapstructure:"description" json:"description"`

	// Dependencies
	Logger string `mapstructure:"logger" json:"logger"`
}

//go:generate go run github.com/invopop/jsonschema -o schema.json -package models AppOptions
