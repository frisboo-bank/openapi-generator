package models

type AppOptions struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Description string `mapstructure:"description"`

	// Dependencies
	Logger string `mapstructure:"logger"`
}
