package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rodrwan/themenu/internal/cqrs/commands"
	"github.com/rodrwan/themenu/internal/utils"
)

type OrderHandler struct {
	commandBus commands.CommandDispatcher
}

func NewOrderHandler(commandBus commands.CommandDispatcher) *OrderHandler {
	return &OrderHandler{
		commandBus: commandBus,
	}
}

// CreateOrder maneja la creación de una nueva orden
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var request struct {
		DishID string `json:"dish_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos"})
		return
	}

	// Obtener el ID del usuario del contexto (asumiendo que viene del middleware de autenticación)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Log temporal para depuración
	log.Printf("userID (from context): %v, tipo: %T", userID, userID)

	var userUUID uuid.UUID
	switch v := userID.(type) {
	case uuid.UUID:
		userUUID = v
	case pgtype.UUID:
		userUUID = utils.FromPgUUID(v)
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor: user_id string inválido"})
			return
		}
		userUUID = parsed
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor: tipo de user_id inesperado"})
		return
	}

	// Convertir los IDs a UUID
	dishUUID, err := uuid.Parse(request.DishID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de plato inválido"})
		return
	}

	// Crear y ejecutar el comando
	cmd := &commands.CreateOrderCommand{
		UserID:  userUUID,
		DishID:  dishUUID,
		Queries: nil, // Se establecerá en el handler
	}

	if err := h.commandBus.Dispatch(cmd); err != nil {
		switch err {
		case commands.ErrOrderExists:
			c.JSON(http.StatusConflict, gin.H{"error": "Ya tienes una orden activa"})
		case commands.ErrDishNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Plato no encontrado"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la orden"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Orden creada exitosamente"})
}
