package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Reader interface {
	GetDish(ctx context.Context, id pgtype.UUID) (Dish, error)
	GetDishByName(ctx context.Context, name string) (Dish, error)
	GetNotificationsByUserId(ctx context.Context, userID pgtype.UUID) ([]Notification, error)
	GetOrder(ctx context.Context, id pgtype.UUID) (Order, error)
	GetOrdersByDishId(ctx context.Context, dishID pgtype.UUID) ([]Order, error)
	GetOrdersByStatus(ctx context.Context, status string) ([]Order, error)
	GetOrdersByUserId(ctx context.Context, userID pgtype.UUID) ([]Order, error)
	GetPermissions(ctx context.Context) ([]Permission, error)
	GetRolePermissions(ctx context.Context) ([]RolePermission, error)
	GetRoles(ctx context.Context) ([]Role, error)
	GetUser(ctx context.Context, id pgtype.UUID) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserRoles(ctx context.Context) ([]UserRole, error)
}

type Writer interface {
	CreateDish(ctx context.Context, arg CreateDishParams) (Dish, error)
	CreateNotification(ctx context.Context, arg CreateNotificationParams) (Notification, error)
	CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteDish(ctx context.Context, id pgtype.UUID) error
	DeleteUser(ctx context.Context, id pgtype.UUID) error
	UpdateDish(ctx context.Context, arg UpdateDishParams) (Dish, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}
