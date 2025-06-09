package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rodrwan/themenu/internal/database"
	"github.com/rodrwan/themenu/internal/utils"
)

type DishHandler struct {
	db database.Querier
}

func NewDishHandler(db database.Querier) *DishHandler {
	return &DishHandler{
		db: db,
	}
}

// ListDishes maneja la obtención de la lista de platos
func (h *DishHandler) ListDishes(c *gin.Context) {
	dishes, err := h.db.ListDishes(c.Request.Context())
	if err != nil {
		log.Printf("Error al obtener los platos: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los platos"})
		return
	}

	response := make([]gin.H, len(dishes))
	for i, dish := range dishes {
		response[i] = gin.H{
			"id":                utils.FromPgUUID(dish.ID).String(),
			"name":              dish.Name,
			"description":       dish.Description.String,
			"price":             utils.ToFloat64(dish.Price),
			"prep_time_minutes": dish.PrepTimeMinutes,
			"available_on":      dish.AvailableOn.Time,
			"created_at":        dish.CreatedAt.Time,
			"updated_at":        dish.UpdatedAt.Time,
		}
	}

	c.JSON(http.StatusOK, response)
}
