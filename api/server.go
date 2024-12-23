package api

import (
	db "github.com/T-BO0/bank/db/sqlc"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Server serves HTTP requests for our banking service
type Server struct {
	store  db.Store
	router *echo.Echo
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Start(address)
}

// NewServer creates a new HTTP server and setup routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := echo.New()
	router.Validator = &CustomValidator{validator: validator.New()}

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.getListOfAccount)

	router.POST("/transfers", server.createTransfer)
	router.GET("/transfers/:id", server.getTransfer)
	router.GET("/transfers", server.listTransfers)

	router.POST("/users", server.createUser)
	router.GET("/users/*", server.getUser)

	server.router = router
	return server
}
