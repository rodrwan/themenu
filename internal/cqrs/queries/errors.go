package queries

import "errors"

var (
	ErrInvalidQuery = errors.New("consulta inválida")
	ErrMenuNotFound = errors.New("menú no encontrado para la fecha especificada")
)
