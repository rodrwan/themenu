package commands

import (
	"context"

	"github.com/google/uuid"
	"github.com/rodrwan/themenu/internal/cqrs"
	"github.com/rodrwan/themenu/internal/database"
	"github.com/rodrwan/themenu/internal/utils"
)

// UpdateOrderStatusCommand representa el comando para actualizar el estado de una orden
type UpdateOrderStatusCommand struct {
	OrderID  uuid.UUID
	Status   string
	Queries  database.Querier
	EventBus *cqrs.EventBus
}

// Execute implementa la interfaz Command
func (c *UpdateOrderStatusCommand) Execute() error {
	ctx := context.Background()

	// Verificar si la orden existe
	order, err := c.Queries.GetOrder(ctx, utils.ToPgUUID(c.OrderID))
	if err != nil {
		return ErrOrderNotFound
	}

	// Actualizar el estado
	_, err = c.Queries.UpdateOrderStatus(ctx, database.UpdateOrderStatusParams{
		ID:     utils.ToPgUUID(c.OrderID),
		Status: c.Status,
	})
	if err != nil {
		return err
	}

	// Publicar evento de actualización de estado
	c.EventBus.PublishEvent(cqrs.EventOrderStatusUpdated, c.Status, map[string]interface{}{
		"order_id": c.OrderID,
		"user_id":  utils.FromPgUUID(order.UserID),
		"dish_id":  utils.FromPgUUID(order.DishID),
		"status":   c.Status,
	})

	return nil
}

// UpdateOrderStatusHandler maneja el comando UpdateOrderStatus
type UpdateOrderStatusHandler struct {
	db       database.Querier
	eventBus *cqrs.EventBus
}

// NewUpdateOrderStatusHandler crea una nueva instancia del handler
func NewUpdateOrderStatusHandler(db database.Querier, eventBus *cqrs.EventBus) *UpdateOrderStatusHandler {
	return &UpdateOrderStatusHandler{
		db:       db,
		eventBus: eventBus,
	}
}

// Handle implementa la interfaz CommandHandler
func (h *UpdateOrderStatusHandler) Handle(command Command) error {
	cmd, ok := command.(*UpdateOrderStatusCommand)
	if !ok {
		return ErrInvalidCommand
	}
	cmd.Queries = h.db
	cmd.EventBus = h.eventBus
	return cmd.Execute()
}
