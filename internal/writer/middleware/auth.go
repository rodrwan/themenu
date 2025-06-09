package middleware

import (
	"net/http"
	"strings"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rodrwan/themenu/internal/database"
	"github.com/rodrwan/themenu/internal/utils"
)

// AuthMiddleware verifica la autenticación del usuario
func AuthMiddleware(db database.Querier) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el token del header
		authHeader := c.GetHeader("Authorization")
		log.Printf("Middleware: authHeader: %v", authHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token no proporcionado"})
			c.Abort()
			return
		}

		// Verificar el formato del token
		parts := strings.Split(authHeader, " ")
		log.Printf("Middleware: parts del token: %v", parts)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			c.Abort()
			return
		}

		// TODO: Implementar la validación real del token JWT
		// Por ahora, asumimos que el token es el ID del usuario
		userID, err := uuid.Parse(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Log temporal para depuración
		log.Printf("Middleware: userID (parsed): %v, tipo: %T", userID, userID)

		// Verificar que el usuario existe
		user, err := db.GetUser(c.Request.Context(), utils.ToPgUUID(userID))
		if err != nil {
			log.Printf("Middleware: error al obtener usuario: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
			c.Abort()
			return
		}

		// Log temporal para depuración
		log.Printf("Middleware: usuario encontrado: %v", user)

		// Guardar el ID del usuario en el contexto
		c.Set("user_id", user.ID) // Guardar como uuid.UUID
		c.Next()
	}
}
