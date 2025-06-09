package commands

import (
	"context"

	"github.com/google/uuid"
	"github.com/rodrwan/themenu/internal/cqrs"
	"github.com/rodrwan/themenu/internal/database"
	"github.com/rodrwan/themenu/internal/utils"
)

// CreateOrderCommand representa el comando para crear una nueva orden
type CreateOrderCommand struct {
	UserID   uuid.UUID
	DishID   uuid.UUID
	Queries  database.Querier
	EventBus *cqrs.EventBus
}

// Execute implementa la interfaz Command
func (c *CreateOrderCommand) Execute() error {
	ctx := context.Background()

	// Verificar si el usuario ya tiene una orden activa
	userUUID := utils.ToPgUUID(c.UserID)
	orders, err := c.Queries.GetOrdersByUserId(ctx, userUUID)
	if err != nil {
		return err
	}

	for _, order := range orders {
		if order.Status != "served" && order.Status != "cancelled" {
			return ErrOrderExists
		}
	}

	// Verificar si el plato existe
	dishUUID := utils.ToPgUUID(c.DishID)
	_, err = c.Queries.GetDish(ctx, dishUUID)
	if err != nil {
		return ErrDishNotFound
	}

	// Crear la orden
	orderID := uuid.New()
	_, err = c.Queries.CreateOrder(ctx, database.CreateOrderParams{
		ID:     utils.ToPgUUID(orderID),
		UserID: userUUID,
		DishID: dishUUID,
		Status: "received",
	})
	if err != nil {
		return err
	}

	// Publicar evento de orden creada
	c.EventBus.PublishEvent(cqrs.EventOrderCreated, "received", map[string]interface{}{
		"order_id": orderID,
		"user_id":  c.UserID,
		"dish_id":  c.DishID,
		"status":   "received",
	})

	return nil
}

// CreateOrderHandler maneja el comando CreateOrder
type CreateOrderHandler struct {
	db       database.Querier
	eventBus *cqrs.EventBus
}

// NewCreateOrderHandler crea una nueva instancia del handler
func NewCreateOrderHandler(db database.Querier, eventBus *cqrs.EventBus) *CreateOrderHandler {
	return &CreateOrderHandler{
		db:       db,
		eventBus: eventBus,
	}
}

// Handle implementa la interfaz CommandHandler
func (h *CreateOrderHandler) Handle(command Command) error {
	cmd, ok := command.(*CreateOrderCommand)
	if !ok {
		return ErrInvalidCommand
	}
	cmd.Queries = h.db
	cmd.EventBus = h.eventBus
	return cmd.Execute()
}
