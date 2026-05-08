package api

import (
	"fmt"
	db "interview/db/sqlc"
	"interview/db/util"
	token "interview/token"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	Config     util.Config
	Store      db.StoreInterface
	Router     *gin.Engine
	TokenMaker token.Maker
}

func NewServer(config util.Config, store db.StoreInterface) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		Config:     config,
		Store:      store,
		TokenMaker: tokenMaker,
	}

	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	authRoutes := router.Group("/").Use(authMiddleware(server.TokenMaker))
	authRoutes.POST("/account", server.CreateAccount)
	authRoutes.GET("/account/:id", server.GetAccount)
	authRoutes.GET("/accounts", server.ListAccounts)

	authRoutes.POST("/transfer", server.CreateTransfer)

	router.POST("/users", server.CreateUser)
	router.POST("/users/login", server.LoginUser)
	server.Router = router
}

func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}
