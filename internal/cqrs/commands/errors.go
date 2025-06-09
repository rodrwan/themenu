package commands

import "errors"

var (
	ErrInvalidCommand = errors.New("comando inválido")
	ErrOrderExists    = errors.New("el usuario ya tiene una orden activa")
	ErrDishNotFound   = errors.New("plato no encontrado")
	ErrOrderNotFound  = errors.New("order not found")
)
