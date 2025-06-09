package cqrs

// Tipos de eventos del sistema
const (
	// Eventos de Usuario
	EventUserCreated = "UserCreated"
	EventUserUpdated = "UserUpdated"
	EventUserDeleted = "UserDeleted"

	// Eventos de Plato
	EventDishCreated = "DishCreated"
	EventDishUpdated = "DishUpdated"
	EventDishDeleted = "DishDeleted"

	// Eventos de Orden
	EventOrderCreated       = "OrderCreated"
	EventOrderStatusUpdated = "OrderStatusUpdated"
	EventOrderCancelled     = "OrderCancelled"
	EventOrderUpdated       = "OrderUpdated"
	EventOrderDeleted       = "OrderDeleted"

	// Eventos de Notificación
	EventNotificationSent = "NotificationSent"

	// Eventos de Sistema
	EventSystemError = "SystemError"

	// Eventos de token
	EventTokenGenerated = "token.generated"
)

// Payloads de eventos
type (
	// UserEventPayload representa el payload para eventos de usuario
	UserEventPayload struct {
		UserID    string `json:"user_id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Timestamp string `json:"timestamp"`
	}

	// DishEventPayload representa el payload para eventos de plato
	DishEventPayload struct {
		DishID          string  `json:"dish_id"`
		Name            string  `json:"name"`
		Description     string  `json:"description"`
		Price           float64 `json:"price"`
		PrepTimeMinutes int     `json:"prep_time_minutes"`
		AvailableOn     string  `json:"available_on"`
		Timestamp       string  `json:"timestamp"`
	}

	// OrderEventPayload representa el payload para eventos de orden
	OrderEventPayload struct {
		OrderID   string `json:"order_id"`
		UserID    string `json:"user_id"`
		DishID    string `json:"dish_id"`
		Status    string `json:"status"`
		Timestamp string `json:"timestamp"`
	}

	// NotificationEventPayload representa el payload para eventos de notificación
	NotificationEventPayload struct {
		NotificationID string `json:"notification_id"`
		UserID         string `json:"user_id"`
		OrderID        string `json:"order_id"`
		Message        string `json:"message"`
		Timestamp      string `json:"timestamp"`
	}

	// SystemEventPayload representa el payload para eventos del sistema
	SystemEventPayload struct {
		Error     string `json:"error"`
		Component string `json:"component"`
		Timestamp string `json:"timestamp"`
	}

	// TokenEventPayload representa el payload del evento de generación de token
	TokenEventPayload struct {
		UserID    string `json:"user_id"`
		Email     string `json:"email"`
		Timestamp string `json:"timestamp"`
	}
)
