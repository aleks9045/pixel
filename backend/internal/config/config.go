package config

import (
	"os"
	"strconv"
)

type OAuthProviderConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type PixelFieldConfig struct {
	FieldWidth  int
	FieldHeight int
}

type Config struct {
	HomepageURL string

	OAuthProviders map[string]OAuthProviderConfig

	PixelFieldConfig PixelFieldConfig

	JwtSecret string
}

func Load() *Config {
	width, err := strconv.Atoi(os.Getenv("PIXEL_FIELD_WIDTH"))
	if err != nil {
		panic(err)
	}
	height, err := strconv.Atoi(os.Getenv("PIXEL_FIELD_HEIGHT"))
	if err != nil {
		panic(err)
	}
	return &Config{
		HomepageURL: os.Getenv("HOMEPAGE_URL"),
		OAuthProviders: map[string]OAuthProviderConfig{
			"github": {
				ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
				ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
				RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
			},
		},
		PixelFieldConfig: PixelFieldConfig{
			width,
			height,
		},

		JwtSecret: os.Getenv("JWT_SECRET"),
	}
}
