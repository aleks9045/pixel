package provider

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aleks9045/pixels/internal/config"
	"github.com/aleks9045/pixels/internal/model"
	"github.com/patrickmn/go-cache"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GitHubAuthProvider struct {
	Oauth oauth2.Config
}

func NewGitHubAuthProvider(cfg *config.OAuthProviderConfig) *GitHubAuthProvider {
	return &GitHubAuthProvider{
		Oauth: oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Endpoint:     github.Endpoint,
			Scopes:       []string{"user:email", "read:user"},
		},
	}
}

func (g *GitHubAuthProvider) GetAuthUrl(state string) (string, error) {
	err := oauthStateCache.Add(state, nil, cache.DefaultExpiration)
	if err != nil {
		return "", errors.New(err.Error())
	}
	return g.Oauth.AuthCodeURL(state), nil
}

func (g *GitHubAuthProvider) Exchange(state, code string) (*model.OauthUser, error) {
	if _, exists := oauthStateCache.Get(state); !exists {
		return nil, errors.New("Wrong state")
	}
	oauthStateCache.Delete(state)

	token, err := g.Oauth.Exchange(context.Background(), code)
	if err != nil {
		return nil, errors.New("Unable to exchange access token: " + err.Error())
	}

	client := g.Oauth.Client(context.Background(), token)
	resp, err := client.Get(gitHubUserInfoURL)
	if err != nil {
		return nil, errors.New("Unable to retrieve user information: " + err.Error())
	}
	defer resp.Body.Close()
	var oauthUser model.OauthUser
	if err := json.NewDecoder(resp.Body).Decode(&oauthUser); err != nil {
		return nil, errors.New("Unable to parse user information: " + err.Error())
	}

	return &oauthUser, nil
}
