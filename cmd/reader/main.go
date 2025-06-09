package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rodrwan/themenu/internal/cqrs/queries"
	"github.com/rodrwan/themenu/internal/database"
	"github.com/rodrwan/themenu/internal/reader"
)

func main() {
	// Configurar la conexión a la base de datos
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/themenu?sslmode=disable"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}
	defer pool.Close()

	// Crear el cliente de base de datos
	db := database.New(pool)

	// Configurar los buses
	qryBus := queries.NewQueryBus()

	// Registrar los handlers
	qryBus.Register("GetMenu", queries.NewGetMenuHandler(db))
	qryBus.Register("GetUserOrders", queries.NewGetUserOrdersHandler(db))

	// Crear y configurar el servidor
	server := reader.NewServer(qryBus, db)

	// Iniciar el servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Servidor iniciado en el puerto %s", port)
	if err := server.Start(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
