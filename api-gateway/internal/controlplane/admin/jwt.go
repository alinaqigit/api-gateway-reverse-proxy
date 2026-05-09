package admin

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const defaultTokenTTL = time.Hour

type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

type AdminClaims struct {
	AdminID string `json:"admin_id"`
	Email   string `json:"email"`
	Elevated bool  `json:"elevated"`
	jwt.RegisteredClaims
}

func NewJWTManagerFromSecret(secret string) (*JWTManager, error) {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return nil, errors.New("JWT_SECRET is required")
	}

	return &JWTManager{
		secret: []byte(secret),
		ttl:    defaultTokenTTL,
	}, nil
}

func (m *JWTManager) GenerateToken(adminID string, email string, elevated bool) (string, time.Time, error) {
	expiresAt := time.Now().Add(m.ttl)
	claims := AdminClaims{
		AdminID: adminID,
		Email:   email,
		Elevated: elevated,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiresAt, nil
}

func (m *JWTManager) ParseToken(tokenString string) (*AdminClaims, error) {
	parsed, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsed.Claims.(*AdminClaims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func AuthMiddleware(jwtManager *JWTManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "missing authorization header"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "invalid authorization header"})
			return
		}

		claims, err := jwtManager.ParseToken(parts[1])
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		ctx.Set("admin_id", claims.AdminID)
		ctx.Set("admin_email", claims.Email)
		ctx.Set("admin_elevated", claims.Elevated)
		ctx.Next()
	}
}

func RequireElevated() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		value, exists := ctx.Get("admin_elevated")
		if !exists {
			ctx.AbortWithStatusJSON(403, gin.H{"error": "admin privileges required"})
			return
		}

		elevated, ok := value.(bool)
		if !ok || !elevated {
			ctx.AbortWithStatusJSON(403, gin.H{"error": "admin privileges required"})
			return
		}

		ctx.Next()
	}
}
