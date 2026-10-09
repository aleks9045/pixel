package oauth

import (
	"net/http"

	"github.com/aleks9045/pixels/internal/config"
	"github.com/aleks9045/pixels/internal/repository"
	"github.com/aleks9045/pixels/internal/service/auth"
	"github.com/aleks9045/pixels/internal/service/oauth/provider"
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type OauthService struct {
	cfg            *config.Config
	authService    *auth.AuthService
	userRepository *repository.UserRepository
	providers      map[string]provider.OauthProvider
}

func NewOauthService(
	cfg *config.Config,
	authService *auth.AuthService,
	userRepository *repository.UserRepository,
	providers map[string]provider.OauthProvider) *OauthService {
	return &OauthService{
		cfg:            cfg,
		userRepository: userRepository,
		authService:    authService,
		providers:      providers,
	}
}

func (s *OauthService) Login(c *gin.Context) {
	oauth := c.Param("oauth")

	provider, ok := s.providers[oauth]
	if !ok {
		c.String(http.StatusNotFound, "Unknown provider")
		return
	}

	state, err := GenerateRandomState()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	redirectUrl, err := provider.GetAuthUrl(state)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	c.Redirect(http.StatusFound, redirectUrl)
}

func (s *OauthService) Callback(c *gin.Context) {
	oauth := c.Param("oauth")

	provider, ok := s.providers[oauth]
	if !ok {
		c.String(http.StatusNotFound, "Unknown provider")
		return
	}

	state := c.Query("state")
	if state == "" {
		c.String(http.StatusBadRequest, "State not provided")
		return
	}

	code := c.Query("code")
	if code == "" {
		c.String(http.StatusBadRequest, "Authorization code not provided")
		return
	}

	oauthUser, err := provider.Exchange(state, code)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	nanoId, err := gonanoid.New()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	userID, err := s.userRepository.CreateIfNotExists(oauthUser, nanoId)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	s.authService.Authenticate(c, userID)

	c.Redirect(http.StatusFound, s.cfg.HomepageURL)
}
