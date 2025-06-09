-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (id, name, email)
VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET name = $2, email = $3
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetDish :one
SELECT * FROM dishes
WHERE id = $1 LIMIT 1;

-- name: GetDishByName :one
SELECT * FROM dishes
WHERE name = $1 LIMIT 1;

-- name: CreateDish :one
INSERT INTO dishes (
    id,
    name,
    description,
    price,
    prep_time_minutes,
    available_on
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: UpdateDish :one
UPDATE dishes
SET
    name = $2,
    description = $3,
    price = $4,
    prep_time_minutes = $5,
    available_on = $6,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteDish :exec
DELETE FROM dishes
WHERE id = $1;

-- name: ListDishes :many
SELECT id, name, description, price, prep_time_minutes, available_on, created_at, updated_at FROM dishes
ORDER BY created_at DESC;

-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- name: CreateOrder :one
INSERT INTO orders (id, user_id, dish_id, status)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetOrdersByUserId :many
SELECT
    o.*,
    d.name as dish_name,
    d.description as dish_description,
    d.price as dish_price
FROM orders o
JOIN dishes d ON o.dish_id = d.id
WHERE o.user_id = $1;

-- name: GetOrdersByDishId :many
SELECT * FROM orders
WHERE dish_id = $1;

-- name: GetOrdersByStatus :many
SELECT * FROM orders
WHERE status = $1;

-- name: GetNotificationsByUserId :many
SELECT * FROM notifications
WHERE user_id = $1;

-- name: CreateNotification :one
INSERT INTO notifications (id, user_id, order_id, message)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetRoles :many
SELECT * FROM roles;

-- name: GetPermissions :many
SELECT * FROM permissions;

-- name: GetRolePermissions :many
SELECT * FROM role_permissions;

-- name: GetUserRoles :many
SELECT * FROM user_roles;

-- name: GetDishesByDate :many
SELECT * FROM dishes
WHERE available_on = $1
ORDER BY name;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $2,
    updated_at = now()
WHERE id = $1
RETURNING id, user_id, dish_id, status, created_at, updated_at;