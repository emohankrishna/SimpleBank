package api

import (
	"fmt"

	db "github.com/emohankrishna/Simplebank/db/sqlc"
	"github.com/emohankrishna/Simplebank/db/util"
	"github.com/emohankrishna/Simplebank/token"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// New Server Instance
func NewServer(config util.Config, store db.Store) (*Server, error) {

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validateCurrency)
	}
	maker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: maker,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.CreateUser)
	router.POST("/users/login", server.loginUser)
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/accounts", server.CreateAccount)
	authRoutes.GET("/accounts/:id", server.GetAccount)
	authRoutes.GET("/accounts", server.ListAccounts)
	authRoutes.POST("/transfers", server.CreateTransfer)

	server.router = router
}

// Start runs the HTTP server on a specif address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
