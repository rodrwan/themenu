package queries

import "sync"

// QueryBus implementa el bus de consultas
type QueryBus struct {
	handlers map[string]QueryHandler
	mu       sync.RWMutex
}

// NewQueryBus crea una nueva instancia del QueryBus
func NewQueryBus() *QueryBus {
	return &QueryBus{
		handlers: make(map[string]QueryHandler),
	}
}

// Register registra un handler para un tipo de consulta
func (b *QueryBus) Register(queryType string, handler QueryHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[queryType] = handler
}

// Dispatch envía una consulta a su handler correspondiente
func (b *QueryBus) Dispatch(query Query) (interface{}, error) {
	b.mu.RLock()
	handler, exists := b.handlers[getQueryType(query)]
	b.mu.RUnlock()

	if !exists {
		return nil, ErrInvalidQuery
	}

	return handler.Handle(query)
}

// getQueryType retorna el tipo de consulta
func getQueryType(query Query) string {
	switch query.(type) {
	case *GetMenuQuery:
		return "GetMenu"
	case *GetUserOrdersQuery:
		return "GetUserOrders"
	default:
		return "Unknown"
	}
}
