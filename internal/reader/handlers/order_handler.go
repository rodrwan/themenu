package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rodrwan/themenu/internal/cqrs/queries"
	"github.com/rodrwan/themenu/internal/utils"
)

type OrderHandler struct {
	queryBus queries.QueryDispatcher
}

func NewOrderHandler(queryBus queries.QueryDispatcher) *OrderHandler {
	return &OrderHandler{
		queryBus: queryBus,
	}
}

// GetMenu maneja la obtención del menú del día
func (h *OrderHandler) GetMenu(c *gin.Context) {
	// Obtener la fecha del query parameter o usar la fecha actual
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido"})
		return
	}

	// Crear y ejecutar la consulta
	query := &queries.GetMenuQuery{
		Date:    date,
		Queries: nil, // Se establecerá en el handler
	}

	result, err := h.queryBus.Dispatch(query)
	if err != nil {
		switch err {
		case queries.ErrMenuNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "No hay menú disponible para esta fecha"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el menú"})
		}
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetUserOrders maneja la obtención de las órdenes del usuario
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	// Obtener el ID del usuario del contexto
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

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

	// Crear y ejecutar la consulta
	query := &queries.GetUserOrdersQuery{
		UserID: userUUID,
	}

	result, err := h.queryBus.Dispatch(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener las órdenes"})
		return
	}

	c.JSON(http.StatusOK, result)
}
