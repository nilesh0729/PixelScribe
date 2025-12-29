package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/nilesh0729/PixelScribe/Result"
	"github.com/nilesh0729/PixelScribe/token"
	"github.com/nilesh0729/PixelScribe/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	TokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		TokenMaker: tokenMaker,
	}
	router := gin.Default()
	router.Use(corsMiddleware())

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.TokenMaker))

	authRoutes.GET("/users/:username", server.getUser)

	authRoutes.POST("/tts/generate", server.generateTTS)
	authRoutes.POST("/dictations", server.createDictation)
	authRoutes.GET("/dictations", server.listDictations)
	authRoutes.DELETE("/dictations/:id", server.deleteDictation)

	authRoutes.POST("/attempts", server.submitAttempt)
	authRoutes.GET("/attempts", server.listAttempts)
    authRoutes.GET("/attempts/:id", server.getAttempt)

	authRoutes.GET("/settings", server.getSettings)
	authRoutes.PUT("/settings", server.updateSettings)

	authRoutes.GET("/performance", server.listPerformance)
	authRoutes.GET("/performance/recent", server.getOverallPerformance)

	server.router = router
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

