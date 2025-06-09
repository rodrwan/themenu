package queries

import (
	"context"
	"time"

	"github.com/rodrwan/themenu/internal/database"
	"github.com/rodrwan/themenu/internal/utils"
)

// MenuItem representa un plato en el menú
type MenuItem struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Price           float64   `json:"price"`
	PrepTimeMinutes int       `json:"prep_time_minutes"`
	AvailableOn     time.Time `json:"available_on"`
}

// GetMenuQuery representa la consulta para obtener el menú del día
type GetMenuQuery struct {
	Date    time.Time
	Queries database.Querier
}

func (q *GetMenuQuery) Execute() (interface{}, error) {
	ctx := context.Background()

	// Obtener los platos disponibles para la fecha especificada
	dishes, err := q.Queries.GetDishesByDate(ctx, utils.ToPgDate(q.Date))
	if err != nil {
		return nil, err
	}

	if len(dishes) == 0 {
		return nil, ErrMenuNotFound
	}

	// Convertir los platos al formato de respuesta
	menuItems := make([]MenuItem, len(dishes))
	for i, dish := range dishes {
		menuItems[i] = MenuItem{
			ID:              utils.FromPgUUID(dish.ID).String(),
			Name:            dish.Name,
			Description:     dish.Description.String,
			Price:           utils.ToFloat64(dish.Price),
			PrepTimeMinutes: int(dish.PrepTimeMinutes),
			AvailableOn:     dish.AvailableOn.Time,
		}
	}

	return menuItems, nil
}

// GetMenuHandler maneja la consulta del menú
type GetMenuHandler struct {
	queries database.Querier
}

func NewGetMenuHandler(queries database.Querier) *GetMenuHandler {
	return &GetMenuHandler{
		queries: queries,
	}
}

func (h *GetMenuHandler) Handle(query Query) (interface{}, error) {
	q, ok := query.(*GetMenuQuery)
	if !ok {
		return nil, ErrInvalidQuery
	}
	q.Queries = h.queries
	return q.Execute()
}
