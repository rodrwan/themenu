package reader

import (
	"github.com/gin-gonic/gin"
	"github.com/rodrwan/themenu/internal/cqrs/queries"
	"github.com/rodrwan/themenu/internal/database"
	"github.com/rodrwan/themenu/internal/reader/handlers"
	"github.com/rodrwan/themenu/internal/reader/middleware"
)

// Server representa el servidor HTTP
type Server struct {
	router   *gin.Engine
	queryBus queries.QueryDispatcher
	db       database.Querier
}

// NewServer crea una nueva instancia del servidor
func NewServer(queryBus queries.QueryDispatcher, db database.Querier) *Server {
	server := &Server{
		router:   gin.Default(),
		queryBus: queryBus,
		db:       db,
	}

	server.setupRoutes()
	return server
}

// setupRoutes configura las rutas del servidor
func (s *Server) setupRoutes() {
	// Middlewares globales
	s.router.Use(middleware.LoggerMiddleware())

	// Aplicar middleware de autenticación para el resto de rutas
	s.router.Use(middleware.AuthMiddleware(s.db))

	orderHandler := handlers.NewOrderHandler(s.queryBus)
	// Rutas protegidas
	menu := s.router.Group("/menu")
	{
		menu.GET("", orderHandler.GetMenu)
	}
	// Rutas de órdenes
	orders := s.router.Group("/orders")
	{
		orders.GET("", orderHandler.GetUserOrders)
	}
	// Rutas de platos
	dishHandler := handlers.NewDishHandler(s.db)
	dishes := s.router.Group("/dishes")
	{
		dishes.GET("", dishHandler.ListDishes)
	}
}

// Start inicia el servidor
func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}
