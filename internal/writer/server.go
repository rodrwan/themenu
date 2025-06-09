package writer

import (
	"github.com/gin-gonic/gin"
	"github.com/rodrwan/themenu/internal/cqrs"
	"github.com/rodrwan/themenu/internal/cqrs/commands"
	"github.com/rodrwan/themenu/internal/database"
	"github.com/rodrwan/themenu/internal/writer/handlers"
	"github.com/rodrwan/themenu/internal/writer/middleware"
)

// Server representa el servidor HTTP
type Server struct {
	router     *gin.Engine
	commandBus commands.CommandDispatcher
	db         database.Querier
	eventBus   *cqrs.EventBus
}

// NewServer crea una nueva instancia del servidor
func NewServer(commandBus commands.CommandDispatcher, db database.Querier, eventBus *cqrs.EventBus) *Server {
	server := &Server{
		router:     gin.Default(),
		commandBus: commandBus,
		db:         db,
		eventBus:   eventBus,
	}

	server.setupRoutes()
	return server
}

// setupRoutes configura las rutas del servidor
func (s *Server) setupRoutes() {
	// Middlewares globales
	s.router.Use(middleware.LoggerMiddleware())

	// Rutas públicas (sin autenticación)
	userHandler := handlers.NewUserHandler(s.db, s.eventBus)
	s.router.POST("/users", userHandler.CreateUser)
	s.router.POST("/users/token", userHandler.GenerateToken)

	// Aplicar middleware de autenticación para el resto de rutas
	s.router.Use(middleware.AuthMiddleware(s.db))

	// Rutas protegidas
	orderHandler := handlers.NewOrderHandler(s.commandBus)
	orders := s.router.Group("/orders")
	{
		orders.POST("", orderHandler.CreateOrder)
		orders.PATCH("/:id/status", orderHandler.UpdateOrderStatus)
	}

	// Rutas de usuario
	users := s.router.Group("/users")
	{
		users.PATCH("/:id", userHandler.UpdateUser)
	}

	// Rutas de platos
	dishHandler := handlers.NewDishHandler(s.db, s.eventBus)
	dishes := s.router.Group("/dishes")
	{
		dishes.POST("", dishHandler.CreateDish)
		dishes.PUT("/:id", dishHandler.UpdateDish)
		dishes.DELETE("/:id", dishHandler.DeleteDish)
	}
}

// Start inicia el servidor
func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}
