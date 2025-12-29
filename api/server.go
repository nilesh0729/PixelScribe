package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/nilesh0729/PixelScribe/Result"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.GET("/users/:username", server.getUser)

	router.POST("/dictations", server.createDictation)
	router.GET("/dictations", server.listDictations)
	router.DELETE("/dictations/:id", server.deleteDictation)

	router.POST("/attempts", server.submitAttempt)
	router.GET("/attempts", server.listAttempts)

	router.GET("/settings", server.getSettings)
	router.PUT("/settings", server.updateSettings)

	router.GET("/performance", server.listPerformance)
	router.GET("/performance/recent", server.getOverallPerformance)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

