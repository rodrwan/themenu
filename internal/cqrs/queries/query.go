package queries

// Query representa una operación de lectura que no modifica el estado
type Query interface {
	Execute() (interface{}, error)
}

// QueryHandler maneja la ejecución de una consulta específica
type QueryHandler interface {
	Handle(query Query) (interface{}, error)
}

// QueryDispatcher es el bus de consultas que distribuye las consultas a sus handlers
type QueryDispatcher interface {
	Dispatch(query Query) (interface{}, error)
	Register(queryType string, handler QueryHandler)
}
