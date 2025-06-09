package web

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rodrwan/themenu/internal/cqrs"
	"github.com/rodrwan/themenu/internal/web/templates"
)

type Server struct {
	app       *fiber.App
	eventBus  *cqrs.EventBus
	apiClient APIClient
}

type APIClient interface {
	GetOrders() ([]Order, error)
	UpdateOrderStatus(orderID, status string) error
}

func NewServer(eventBus *cqrs.EventBus, apiClient APIClient) *Server {
	app := fiber.New()

	app.Use(cors.New())

	// Configurar archivos estáticos
	app.Static("/static", "./internal/web/static")

	server := &Server{
		app:       app,
		eventBus:  eventBus,
		apiClient: apiClient,
	}

	// Rutas
	app.Get("/", server.handleDashboard)
	app.Get("/events", server.handleSSE)
	app.Get("/orders", server.handleOrders)
	app.Patch("/orders/:id/status", server.handleUpdateOrderStatus)

	return server
}

func (s *Server) Start(addr string) error {
	return s.app.Listen(addr)
}

func (s *Server) handleDashboard(c *fiber.Ctx) error {
	// Por ahora, enviamos una lista vacía de eventos
	component := templates.Dashboard([]cqrs.Event{})

	var buf bytes.Buffer
	if err := component.Render(c.Context(), &buf); err != nil {
		return err
	}
	c.Type("html")
	return c.SendString(buf.String())
}

func (s *Server) handleSSE(c *fiber.Ctx) error {
	log.Printf("[SSE] Nueva conexión establecida")

	// Configurar headers SSE
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	// Crear un canal para los eventos
	eventChan := s.eventBus.Subscribe("*")
	ticker := time.NewTicker(10 * time.Second)

	// Limpiar recursos cuando se cierre la conexión
	defer func() {
		ticker.Stop()
		log.Printf("[SSE] Conexión cerrada, recursos liberados")
	}()

	// Configurar el writer para streaming
	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		// Enviar un mensaje inicial
		log.Printf("[SSE] Enviando mensaje inicial")
		w.WriteString("data: connected\n\n")
		w.Flush()

		for {
			select {
			case event, ok := <-eventChan:
				if !ok {
					log.Printf("[SSE] Canal de eventos cerrado")
					return
				}

				log.Printf("[SSE] Evento recibido: %s", event.Type)
				data, err := json.Marshal(event)
				if err != nil {
					log.Printf("[SSE] Error al serializar evento: %v", err)
					continue
				}

				w.WriteString(fmt.Sprintf("data: %s\n\n", string(data)))
				w.Flush()

			case <-ticker.C:
				log.Printf("[SSE] Enviando ping")
				w.WriteString("data: ping\n\n")
				w.Flush()
			}
		}
	})

	return nil
}

func (s *Server) handleOrders(c *fiber.Ctx) error {
	// Get order from api service
	orders, err := s.apiClient.GetOrders()
	if err != nil {
		log.Printf("Error getting orders: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get orders"})
	}
	return c.JSON(orders)
}

func (s *Server) handleUpdateOrderStatus(c *fiber.Ctx) error {
	orderID := c.Params("id")
	var body struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Update order status in api service
	err := s.apiClient.UpdateOrderStatus(orderID, body.Status)
	if err != nil {
		log.Printf("Error updating order status: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to update order status",
			"message": "Error updating order status",
		})
	}

	return c.JSON(fiber.Map{"message": fmt.Sprintf("Order %s updated to %s", orderID, body.Status)})
}
