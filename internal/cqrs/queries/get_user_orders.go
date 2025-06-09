package queries

import (
	"context"

	"github.com/google/uuid"
	"github.com/rodrwan/themenu/internal/database"
	"github.com/rodrwan/themenu/internal/utils"
)

// GetUserOrdersQuery representa la consulta para obtener las órdenes de un usuario
type GetUserOrdersQuery struct {
	UserID  uuid.UUID
	Queries database.Querier
}

// Execute implementa la interfaz Query
func (q *GetUserOrdersQuery) Execute() (interface{}, error) {
	ctx := context.Background()

	// Obtener las órdenes del usuario
	orders, err := q.Queries.GetOrdersByUserId(ctx, utils.ToPgUUID(q.UserID))
	if err != nil {
		return nil, err
	}

	// Convertir las órdenes a un formato más amigable
	var result []map[string]interface{}
	for _, order := range orders {
		result = append(result, map[string]interface{}{
			"id":               utils.FromPgUUID(order.ID).String(),
			"user_id":          utils.FromPgUUID(order.UserID).String(),
			"dish_id":          utils.FromPgUUID(order.DishID).String(),
			"dish_name":        order.DishName,
			"dish_description": order.DishDescription,
			"dish_price":       order.DishPrice,
			"status":           order.Status,
			"created_at":       order.CreatedAt.Time,
			"updated_at":       order.UpdatedAt.Time,
		})
	}

	return result, nil
}

// GetUserOrdersHandler maneja la consulta GetUserOrders
type GetUserOrdersHandler struct {
	db database.Querier
}

// NewGetUserOrdersHandler crea una nueva instancia del handler
func NewGetUserOrdersHandler(db database.Querier) *GetUserOrdersHandler {
	return &GetUserOrdersHandler{
		db: db,
	}
}

// Handle implementa la interfaz QueryHandler
func (h *GetUserOrdersHandler) Handle(query Query) (interface{}, error) {
	q, ok := query.(*GetUserOrdersQuery)
	if !ok {
		return nil, ErrInvalidQuery
	}
	q.Queries = h.db
	return q.Execute()
}
