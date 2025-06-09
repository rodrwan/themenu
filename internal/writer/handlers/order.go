package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rodrwan/themenu/internal/cqrs/commands"
)

// UpdateOrderStatus actualiza el estado de una orden
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Printf("Error parsing order ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de orden inválido"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=received confirmed preparing served cancelled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := &commands.UpdateOrderStatusCommand{
		OrderID: orderID,
		Status:  req.Status,
	}

	if err := h.commandBus.Dispatch(cmd); err != nil {
		log.Printf("Error dispatching command: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Estado de la orden actualizado"})
}
