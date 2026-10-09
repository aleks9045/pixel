package provider

import (
	"log"
	"time"

	"github.com/aleks9045/pixels/internal/config"
	"github.com/aleks9045/pixels/internal/model"
	"github.com/patrickmn/go-cache"
)

const (
	gitHubUserInfoURL = "https://api.github.com/user"
)

var (
	oauthStateCache = cache.New(10*time.Minute, 20*time.Minute)
)

type OauthProvider interface {
	GetAuthUrl(state string) (string, error)
	Exchange(state, code string) (*model.OauthUser, error)
}

type OauthProviderFactory func(cfg config.OAuthProviderConfig) OauthProvider

var OauthProviderRegistry = map[string]OauthProviderFactory{
	"github": func(cfg config.OAuthProviderConfig) OauthProvider {
		return NewGitHubAuthProvider(&cfg)
	},
}

func NewProvidersFromConfig(cfgs map[string]config.OAuthProviderConfig) map[string]OauthProvider {

	providers := make(map[string]OauthProvider, len(cfgs))

	for name, cfg := range cfgs {

		factory, ok := OauthProviderRegistry[name]
		if !ok {
			log.Fatal("provider not found: " + name)
		}

		providers[name] = factory(cfg)
	}
	return providers
}
