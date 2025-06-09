# TheMenu - Sistema de Gestión de Restaurante

## Descripción
TheMenu es un sistema backend para la gestión de restaurantes que permite manejar platos, órdenes y usuarios. Implementa patrones modernos de arquitectura como CQRS y utiliza una base de datos PostgreSQL con consultas generadas automáticamente mediante sqlc. Incluye una interfaz web para monitorear eventos en tiempo real.

## Requisitos Previos
- Go 1.24 o superior
- Docker y Docker Compose
- PostgreSQL 15
- Redis 7

## Estructura del Proyecto

```
.
├── cmd/                    # Puntos de entrada de la aplicación
│   ├── reader/            # Servicio de lectura
│   ├── web/               # Servicio web para visualizar eventos
│   └── writer/            # Servicio de escritura
├── internal/              # Código interno de la aplicación
│   ├── cqrs/             # Implementación del patrón CQRS
│   │   ├── commands/     # Comandos para modificar datos
│   │   ├── queries/      # Consultas para leer datos
│   │   └── events/       # Definición de eventos
│   ├── database/         # Capa de base de datos
│   │   ├── models/       # Modelos generados por sqlc
│   │   └── queries/      # Consultas generadas por sqlc
│   ├── reader/           # Servicio de lectura
│   │   ├── handlers/     # Manejadores de consultas
│   │   └── server.go     # Configuración del servidor de lectura
│   ├── utils/            # Utilidades compartidas
│   │   ├── date.go       # Utilidades para fechas
│   │   ├── numeric.go    # Utilidades para números
│   │   ├── text.go       # Utilidades para texto
│   │   └── uuid.go       # Utilidades para UUIDs
│   ├── web/              # Servicio web principal
│   │   ├── handlers/     # Manejadores HTTP
│   │   └── server.go     # Configuración del servidor web
│   └── writer/           # Servicio de escritura
│       ├── handlers/     # Manejadores de comandos
│       └── server.go     # Configuración del servidor de escritura
├── schema.sql            # Esquema de la base de datos
├── queries.sql           # Consultas SQL para sqlc
├── seed.sql             # Datos iniciales
├── sqlc.yaml            # Configuración de sqlc
├── docker-compose.yml   # Configuración de Docker Compose
├── Dockerfile.web       # Dockerfile para el servicio web
├── Dockerfile.reader    # Dockerfile para el servicio de lectura
├── Dockerfile.writer    # Dockerfile para el servicio de escritura
└── Makefile             # Comandos de automatización
```

## Configuración y Ejecución

### 1. Configuración del Entorno
```bash
# Clonar el repositorio
git clone https://github.com/rodrwan/themenu.git
cd themenu

# Configurar variables de entorno
cp .env.example .env
```

### 2. Iniciar Servicios con Docker
```bash
docker-compose up -d
```

### 3. Generar Código SQL y Templ
```bash
make generate
```

### 4. Ejecutar la Aplicación
```bash
# Iniciar todos los servicios
make run

# O iniciar servicios individualmente
make run-web    # Inicia el servicio web
make run-reader # Inicia el servicio de lectura
make run-writer # Inicia el servicio de escritura
```

### 5. Acceder a la Interfaz Web
Una vez que los servicios estén en ejecución, puedes acceder a la interfaz web de monitoreo en:
```
http://localhost:8080
```

## Componentes Principales

### API
- Endpoints RESTful para gestionar platos, órdenes y usuarios
- Autenticación mediante tokens JWT
- Validación de datos con go-validator
- Documentación OpenAPI/Swagger
- Eventos publicados a través del event bus

### Interfaz Web
- Panel de monitoreo en tiempo real de eventos
- Visualización de eventos por tipo y timestamp
- Filtrado y búsqueda de eventos
- Interfaz responsive y moderna
- Actualización en tiempo real mediante WebSocket

### Base de Datos
- PostgreSQL como base de datos principal
- Esquema con tablas para:
  - Platos (dishes)
  - Órdenes (orders)
  - Usuarios (users)
  - Notificaciones (notifications)
- Consultas generadas automáticamente por sqlc
- Migraciones automáticas

### Eventos
- Sistema de eventos para notificar cambios en el sistema
- Eventos principales:
  - DishCreated
  - DishUpdated
  - OrderCreated
  - OrderStatusChanged
  - UserCreated
- Integración con Redis para el event bus

### Utilidades
- Conversiones entre tipos de Go y PostgreSQL
- Manejo de UUIDs, fechas, números y texto
- Funciones auxiliares para la aplicación
- Helpers para validación y sanitización

## Endpoints Principales

### Gestión de Platos
- `POST /api/v1/dishes` - Crear plato
- `PUT /api/v1/dishes/:id` - Actualizar plato
- `DELETE /api/v1/dishes/:id` - Eliminar plato
- `GET /api/v1/dishes` - Listar platos
- `GET /api/v1/dishes/:id` - Obtener plato por ID

### Gestión de Órdenes
- `POST /api/v1/orders` - Crear orden
- `PATCH /api/v1/orders/:id/status` - Actualizar estado
- `GET /api/v1/orders` - Listar órdenes
- `GET /api/v1/orders/:id` - Obtener orden por ID

### Gestión de Usuarios
- `POST /api/v1/users` - Crear usuario
- `PUT /api/v1/users/:id` - Actualizar usuario
- `DELETE /api/v1/users/:id` - Eliminar usuario
- `GET /api/v1/users/:id` - Obtener usuario por ID

## Contribución
1. Fork el repositorio
2. Crear una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abrir un Pull Request

## Licencia
Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para más detalles. 