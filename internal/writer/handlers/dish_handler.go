package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rodrwan/themenu/internal/cqrs"
	"github.com/rodrwan/themenu/internal/database"
	"github.com/rodrwan/themenu/internal/utils"
)

type DishHandler struct {
	db       database.Querier
	eventBus *cqrs.EventBus
}

func NewDishHandler(db database.Querier, eventBus *cqrs.EventBus) *DishHandler {
	return &DishHandler{
		db:       db,
		eventBus: eventBus,
	}
}

// CreateDish maneja la creación de un nuevo plato
func (h *DishHandler) CreateDish(c *gin.Context) {
	var request struct {
		Name            string    `json:"name" binding:"required"`
		Description     string    `json:"description"`
		Price           float64   `json:"price" binding:"required"`
		PrepTimeMinutes int       `json:"prep_time_minutes" binding:"required"`
		AvailableOn     time.Time `json:"available_on" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos"})
		return
	}

	// Crear el plato
	dishID := uuid.New()
	dish, err := h.db.CreateDish(c.Request.Context(), database.CreateDishParams{
		ID:              utils.ToPgUUID(dishID),
		Name:            request.Name,
		Description:     utils.ToPgText(request.Description),
		Price:           utils.ToPgNumeric(request.Price),
		PrepTimeMinutes: int32(request.PrepTimeMinutes),
		AvailableOn:     utils.ToPgDate(request.AvailableOn),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el plato"})
		return
	}

	// Publicar evento de plato creado
	h.eventBus.PublishEvent(cqrs.EventDishCreated, "success", cqrs.DishEventPayload{
		DishID:          dishID.String(),
		Name:            dish.Name,
		Description:     dish.Description.String,
		Price:           utils.ToFloat64(dish.Price),
		PrepTimeMinutes: int(dish.PrepTimeMinutes),
		AvailableOn:     dish.AvailableOn.Time.Format(time.RFC3339),
		Timestamp:       time.Now().Format(time.RFC3339),
	})

	c.JSON(http.StatusCreated, gin.H{
		"id":                utils.FromPgUUID(dish.ID).String(),
		"name":              dish.Name,
		"description":       dish.Description.String,
		"price":             utils.ToFloat64(dish.Price),
		"prep_time_minutes": dish.PrepTimeMinutes,
		"available_on":      dish.AvailableOn.Time,
	})
}

// UpdateDish maneja la actualización de un plato
func (h *DishHandler) UpdateDish(c *gin.Context) {
	dishID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de plato inválido"})
		return
	}

	var request struct {
		Name            string    `json:"name" binding:"required"`
		Description     string    `json:"description"`
		Price           float64   `json:"price" binding:"required"`
		PrepTimeMinutes int       `json:"prep_time_minutes" binding:"required"`
		AvailableOn     time.Time `json:"available_on" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos"})
		return
	}

	// Actualizar el plato
	dish, err := h.db.UpdateDish(c.Request.Context(), database.UpdateDishParams{
		ID:              utils.ToPgUUID(dishID),
		Name:            request.Name,
		Description:     utils.ToPgText(request.Description),
		Price:           utils.ToPgNumeric(request.Price),
		PrepTimeMinutes: int32(request.PrepTimeMinutes),
		AvailableOn:     utils.ToPgDate(request.AvailableOn),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el plato"})
		return
	}

	// Publicar evento de plato actualizado
	h.eventBus.PublishEvent(cqrs.EventDishUpdated, "success", cqrs.DishEventPayload{
		DishID:          dishID.String(),
		Name:            dish.Name,
		Description:     dish.Description.String,
		Price:           utils.ToFloat64(dish.Price),
		PrepTimeMinutes: int(dish.PrepTimeMinutes),
		AvailableOn:     dish.AvailableOn.Time.Format(time.RFC3339),
		Timestamp:       time.Now().Format(time.RFC3339),
	})

	c.JSON(http.StatusOK, gin.H{
		"id":                utils.FromPgUUID(dish.ID).String(),
		"name":              dish.Name,
		"description":       dish.Description.String,
		"price":             utils.ToFloat64(dish.Price),
		"prep_time_minutes": dish.PrepTimeMinutes,
		"available_on":      dish.AvailableOn.Time,
	})
}

// DeleteDish maneja la eliminación de un plato
func (h *DishHandler) DeleteDish(c *gin.Context) {
	dishID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de plato inválido"})
		return
	}

	// Obtener el plato antes de eliminarlo para el evento
	dish, err := h.db.GetDish(c.Request.Context(), utils.ToPgUUID(dishID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plato no encontrado"})
		return
	}

	// Eliminar el plato
	if err := h.db.DeleteDish(c.Request.Context(), utils.ToPgUUID(dishID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el plato"})
		return
	}

	// Publicar evento de plato eliminado
	h.eventBus.PublishEvent(cqrs.EventDishDeleted, "success", cqrs.DishEventPayload{
		DishID:          dishID.String(),
		Name:            dish.Name,
		Description:     dish.Description.String,
		Price:           utils.ToFloat64(dish.Price),
		PrepTimeMinutes: int(dish.PrepTimeMinutes),
		AvailableOn:     dish.AvailableOn.Time.Format(time.RFC3339),
		Timestamp:       time.Now().Format(time.RFC3339),
	})

	c.JSON(http.StatusOK, gin.H{"message": "Plato eliminado exitosamente"})
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
