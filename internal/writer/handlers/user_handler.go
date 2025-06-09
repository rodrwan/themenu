package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rodrwan/themenu/internal/cqrs"
	"github.com/rodrwan/themenu/internal/database"
	"github.com/rodrwan/themenu/internal/utils"
)

type UserHandler struct {
	db       database.Querier
	eventBus *cqrs.EventBus
}

func NewUserHandler(db database.Querier, eventBus *cqrs.EventBus) *UserHandler {
	return &UserHandler{
		db:       db,
		eventBus: eventBus,
	}
}

// CreateUser maneja la creación de un nuevo usuario
func (h *UserHandler) CreateUser(c *gin.Context) {
	var request struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos"})
		return
	}

	// Crear el usuario
	userID := uuid.New()
	user, err := h.db.CreateUser(c.Request.Context(), database.CreateUserParams{
		ID:    utils.ToPgUUID(userID),
		Name:  request.Name,
		Email: request.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el usuario"})
		return
	}

	// Publicar evento de usuario creado
	h.eventBus.PublishEvent(cqrs.EventUserCreated, "success", cqrs.UserEventPayload{
		UserID:    userID.String(),
		Name:      user.Name,
		Email:     user.Email,
		Timestamp: time.Now().Format(time.RFC3339),
	})

	c.JSON(http.StatusCreated, gin.H{
		"id":    utils.FromPgUUID(user.ID).String(),
		"name":  user.Name,
		"email": user.Email,
	})
}

// UpdateUser maneja la actualización de un usuario
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
		return
	}

	var request struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos"})
		return
	}

	// Actualizar el usuario
	user, err := h.db.UpdateUser(c.Request.Context(), database.UpdateUserParams{
		ID:    utils.ToPgUUID(userID),
		Name:  request.Name,
		Email: request.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el usuario"})
		return
	}

	// Publicar evento de usuario actualizado
	h.eventBus.PublishEvent(cqrs.EventUserUpdated, "success", cqrs.UserEventPayload{
		UserID:    userID.String(),
		Name:      user.Name,
		Email:     user.Email,
		Timestamp: time.Now().Format(time.RFC3339),
	})

	c.JSON(http.StatusOK, gin.H{
		"id":    utils.FromPgUUID(user.ID).String(),
		"name":  user.Name,
		"email": user.Email,
	})
}

// GenerateToken maneja la generación de un token de acceso
func (h *UserHandler) GenerateToken(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos"})
		return
	}

	// Buscar el usuario por email
	user, err := h.db.GetUserByEmail(c.Request.Context(), request.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Generar un token (en este caso, el ID del usuario)
	token := utils.FromPgUUID(user.ID).String()

	// Publicar evento de token generado
	h.eventBus.PublishEvent(cqrs.EventTokenGenerated, "success", cqrs.TokenEventPayload{
		UserID:    token,
		Email:     user.Email,
		Timestamp: time.Now().Format(time.RFC3339),
	})

	c.JSON(http.StatusOK, gin.H{"token": token})
}
