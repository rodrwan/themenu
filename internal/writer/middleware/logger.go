package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware registra información sobre las peticiones HTTP
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Tiempo de inicio
		start := time.Now()

		// Procesar la petición
		c.Next()

		// Tiempo de finalización
		end := time.Now()
		latency := end.Sub(start)

		// Obtener información de la petición
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path

		// Formatear y registrar el log
		logMessage := fmt.Sprintf("[%s] %s %s %d %v",
			method,
			path,
			clientIP,
			statusCode,
			latency,
		)

		// TODO: Implementar un sistema de logging más robusto
		fmt.Println(logMessage)
	}
}
