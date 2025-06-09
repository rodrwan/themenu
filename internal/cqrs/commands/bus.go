package commands

import "sync"

// CommandBus implementa el bus de comandos
type CommandBus struct {
	handlers map[string]CommandHandler
	mu       sync.RWMutex
}

// NewCommandBus crea una nueva instancia del CommandBus
func NewCommandBus() *CommandBus {
	return &CommandBus{
		handlers: make(map[string]CommandHandler),
	}
}

// Register registra un handler para un tipo de comando
func (b *CommandBus) Register(commandType string, handler CommandHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[commandType] = handler
}

// Dispatch envía un comando a su handler correspondiente
func (b *CommandBus) Dispatch(command Command) error {
	b.mu.RLock()
	handler, exists := b.handlers[getCommandType(command)]
	b.mu.RUnlock()

	if !exists {
		return ErrInvalidCommand
	}

	return handler.Handle(command)
}

// getCommandType retorna el tipo de comando
func getCommandType(command Command) string {
	switch command.(type) {
	case *CreateOrderCommand:
		return "CreateOrder"
	case *UpdateOrderStatusCommand:
		return "UpdateOrderStatus"
	default:
		return "Unknown"
	}
}
