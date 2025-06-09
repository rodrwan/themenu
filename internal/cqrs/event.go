package cqrs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Event representa un evento en el sistema
type Event struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Payload   string    `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
}

// EventBus maneja la distribución de eventos
type EventBus struct {
	redisClient  *redis.Client
	subscribers  map[string][]chan Event
	ctx          context.Context
	cancel       context.CancelFunc
	mu           sync.RWMutex
	backpressure chan struct{}
}

const (
	// BufferSize define el tamaño del buffer para los canales de eventos
	BufferSize = 1000
	// MaxRetries define el número máximo de intentos para enviar un evento
	MaxRetries = 3
	// RetryDelay define el tiempo de espera entre reintentos
	RetryDelay = 100 * time.Millisecond
	// BackpressureLimit define el límite de eventos en el buffer antes de activar backpressure
	BackpressureLimit = 800
	// RedisBufferSize define el tamaño del buffer para Redis
	RedisBufferSize = 1024 * 1024 // 1MB
)

// NewEventBus crea una nueva instancia de EventBus
func NewEventBus() *EventBus {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}

	// Configurar opciones de reconexión y buffer
	opt.MaxRetries = 3
	opt.MinRetryBackoff = 8 * time.Millisecond
	opt.MaxRetryBackoff = 512 * time.Millisecond
	opt.DialTimeout = 5 * time.Second
	opt.ReadTimeout = 3 * time.Second
	opt.WriteTimeout = 3 * time.Second
	opt.PoolSize = 10
	opt.MinIdleConns = 5

	client := redis.NewClient(opt)

	// Crear contexto con cancelación
	ctx, cancel := context.WithCancel(context.Background())

	bus := &EventBus{
		redisClient:  client,
		subscribers:  make(map[string][]chan Event),
		ctx:          ctx,
		cancel:       cancel,
		backpressure: make(chan struct{}, 1),
	}

	// Verificar la conexión
	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Error al conectar con Redis: %v", err)
		panic(err)
	}

	// Iniciar el subscriber de Redis en una goroutine
	go bus.subscribeToRedis(ctx)

	return bus
}

// subscribeToRedis escucha los eventos publicados en Redis
func (b *EventBus) subscribeToRedis(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			pubsub := b.redisClient.Subscribe(ctx, "events")
			ch := pubsub.Channel(
				redis.WithChannelSize(RedisBufferSize),
			)

			for msg := range ch {
				var event Event
				if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
					log.Printf("Error al deserializar evento: %v", err)
					continue
				}

				// Verificar backpressure
				select {
				case b.backpressure <- struct{}{}:
					// Reenviar el evento a los suscriptores locales (SSE) sin publicarlo en Redis
					b.mu.RLock()
					subscribers := b.subscribers["events"]
					b.mu.RUnlock()
					for _, subscriber := range subscribers {
						select {
						case subscriber <- event:
						default:
							log.Printf("Canal lleno, ignorando evento para suscriptor")
						}
					}
					<-b.backpressure
				default:
					log.Printf("Backpressure activado, ignorando evento")
					time.Sleep(100 * time.Millisecond)
				}
			}

			// Si llegamos aquí, la conexión se cerró
			log.Println("Conexión Redis cerrada, intentando reconectar...")
			time.Sleep(time.Second) // Esperar antes de reconectar
		}
	}
}

// Close cierra la conexión con Redis
func (b *EventBus) Close() {
	b.cancel()
	if err := b.redisClient.Close(); err != nil {
		log.Printf("Error al cerrar la conexión Redis: %v", err)
	}
}

// Subscribe registra un nuevo suscriptor para un tipo de evento
func (b *EventBus) Subscribe(eventType string) <-chan Event {
	ch := make(chan Event, 100)
	b.mu.Lock()
	defer b.mu.Unlock()

	// Si el tipo de evento es vacío o '*', suscribirse al canal general 'events'
	channel := "events"
	if eventType != "" && eventType != "*" {
		channel = fmt.Sprintf("events:%s", eventType)
	}

	b.subscribers[channel] = append(b.subscribers[channel], ch)
	return ch
}

// Unsubscribe elimina un suscriptor
func (b *EventBus) Unsubscribe(eventType string, ch chan Event) {
	b.mu.Lock()
	defer b.mu.Unlock()

	subscribers := b.subscribers[eventType]
	for i, subscriber := range subscribers {
		if subscriber == ch {
			// Cerrar el canal
			close(ch)
			// Eliminar el suscriptor
			b.subscribers[eventType] = append(subscribers[:i], subscribers[i+1:]...)
			return
		}
	}
}

// Publish envía un evento a todos los suscriptores del tipo especificado
func (b *EventBus) Publish(event Event) {
	// Publicar en Redis en el canal específico del tipo de evento
	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error al serializar evento: %v", err)
		return
	}

	// redisChannel := "events:" + event.Type
	// log.Printf("Publicando evento en Redis: %s", redisChannel)
	// if err := b.redisClient.Publish(b.ctx, redisChannel, string(eventBytes)).Err(); err != nil {
	// 	log.Printf("Error al publicar evento en Redis: %v", err)
	// }

	// También publicar en el canal general para compatibilidad
	if err := b.redisClient.Publish(b.ctx, "events", string(eventBytes)).Err(); err != nil {
		log.Printf("Error al publicar evento en Redis: %v", err)
	}
}

// sendWithRetry intenta enviar un evento a un canal con reintentos
func (b *EventBus) sendWithRetry(ch chan Event, event Event) {
	for i := 0; i < MaxRetries; i++ {
		select {
		case ch <- event:
			return
		default:
			if i < MaxRetries-1 {
				time.Sleep(RetryDelay)
			} else {
				log.Printf("No se pudo enviar el evento después de %d intentos", MaxRetries)
			}
		}
	}
}

// PublishEvent es un helper para publicar eventos desde el backend
func (b *EventBus) PublishEvent(eventType, status string, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	event := Event{
		ID:        uuid.New().String(),
		Type:      eventType,
		Status:    status,
		Payload:   string(payloadBytes),
		Timestamp: time.Now(),
	}
	b.Publish(event)
	return nil
}
