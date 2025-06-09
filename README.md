# TheMenu - Documentación del Proyecto

## Estructura del Proyecto

```
.
├── cmd/
│   └── api/
│       └── main.go           # Punto de entrada de la API
├── internal/
│   ├── api/
│   │   ├── handlers/         # Manejadores HTTP
│   │   │   ├── dish_handler.go
│   │   │   ├── order_handler.go
│   │   │   └── user_handler.go
│   │   └── server.go         # Configuración del servidor
│   ├── cqrs/                 # Patrón CQRS
│   │   ├── commands/         # Comandos
│   │   ├── queries/          # Consultas
│   │   └── events.go         # Definición de eventos
│   ├── database/            # Capa de base de datos
│   │   ├── models.go        # Modelos generados por sqlc
│   │   └── query.sql.go     # Consultas generadas por sqlc
│   └── utils/               # Utilidades
│       ├── date.go          # Utilidades para fechas
│       ├── numeric.go       # Utilidades para números
│       ├── text.go          # Utilidades para texto
│       └── uuid.go          # Utilidades para UUIDs
├── schema.sql              # Esquema de la base de datos
├── query.sql              # Consultas SQL para sqlc
├── seed.sql              # Datos iniciales
├── sqlc.yaml            # Configuración de sqlc
└── docker-compose.yml   # Configuración de Docker Compose
```

## Componentes Principales

### API
- Endpoints RESTful para gestionar platos, órdenes y usuarios
- Autenticación mediante tokens
- Eventos publicados a través del event bus

### Base de Datos
- PostgreSQL como base de datos principal
- Esquema con tablas para platos, órdenes, usuarios y notificaciones
- Consultas generadas automáticamente por sqlc

### Eventos
- Sistema de eventos para notificar cambios en el sistema
- Eventos para creación, actualización y eliminación de entidades
- Integración con Redis para el event bus

### Utilidades
- Conversiones entre tipos de Go y PostgreSQL
- Manejo de UUIDs, fechas, números y texto
- Funciones auxiliares para la aplicación

## Flujos Principales

### Gestión de Platos
1. Crear plato (POST /dishes)
2. Actualizar plato (PUT /dishes/:id)
3. Eliminar plato (DELETE /dishes/:id)
4. Listar platos (GET /dishes)

### Gestión de Órdenes
1. Crear orden (POST /orders)
2. Actualizar estado (PATCH /orders/:id/status)
3. Listar órdenes de usuario (GET /orders)

### Gestión de Usuarios
1. Crear usuario (POST /users)
2. Actualizar usuario (PUT /users/:id)
3. Eliminar usuario (DELETE /users/:id)

## Notas Importantes
- Los archivos de base de datos (schema.sql, query.sql, seed.sql) están en la raíz del proyecto
- Las consultas SQL se generan automáticamente con sqlc
- Los eventos se publican a través del event bus de Redis
- Las utilidades para tipos de PostgreSQL están en el paquete utils 