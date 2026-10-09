package main

import (
	"log"
	"os"

	"github.com/aleks9045/pixels/internal/api"
	"github.com/aleks9045/pixels/internal/config"
	"github.com/aleks9045/pixels/internal/db"
	"github.com/aleks9045/pixels/internal/middleware"
	"github.com/aleks9045/pixels/internal/repository"
	"github.com/aleks9045/pixels/internal/service/auth"
	"github.com/aleks9045/pixels/internal/service/oauth"
	"github.com/aleks9045/pixels/internal/service/oauth/provider"
	"github.com/aleks9045/pixels/internal/service/pixel"
	"github.com/aleks9045/pixels/internal/service/user"
	"github.com/aleks9045/pixels/internal/service/ws"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		os.Exit(1)
	}
}

func main() {
	loadEnv()

	app := NewApp()

	app.engine.Run()
}

type App struct {
	engine *gin.Engine
}

func NewApp() *App {

	cfg := config.Load()

	db := db.New()

	tokenManager := auth.NewTokenManager(cfg)
	providers := provider.NewProvidersFromConfig(cfg.OAuthProviders)

	userRepository := repository.NewUserRepository(db)
	pixelRepository := repository.NewPixelRepository(db, cfg.PixelFieldConfig)

	authService := auth.NewAuthService(tokenManager)
	oauthService := oauth.NewOauthService(cfg, authService, userRepository, providers)
	userService := user.NewUserService(userRepository)
	pixelService := pixel.NewPixelService(pixelRepository)
	wsService := ws.NewWsService(pixelService)

	authMiddleware := middleware.NewAuthMiddleware(userRepository, tokenManager)

	engine := gin.New()
	engine.SetTrustedProxies([]string{"localhost"})

	engine.Use(gin.Logger())
	engine.Use(middleware.CORSMiddleware())
	engine.Use(gin.Recovery())

	apiGroup := engine.Group("/api")

	api.RegisterOauthRoutes(apiGroup, oauthService)
	api.RegisterAuthGroup(apiGroup, authService, authMiddleware.Middleware())
	api.RegisterUserRoutes(apiGroup, userService, authMiddleware.Middleware())
	api.RegisterPixelRoutes(apiGroup, pixelService, authMiddleware.Middleware())
	api.RegisterWebsocketRoutes(apiGroup, wsService, authMiddleware.Middleware())

	return &App{
		engine: engine,
	}
}
