package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const (
	UserIDContextKey   contextKey = "userId"
	UserRoleContextKey contextKey = "userRole"
)

type JWTClaims struct {
	UserID string `json:"sub"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"success":false,"error":{"code":"UNAUTHORIZED","message":"Missing authorization token"}}`, http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, `{"success":false,"error":{"code":"UNAUTHORIZED","message":"Invalid authorization header format"}}`, http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]
			claims := &JWTClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, `{"success":false,"error":{"code":"UNAUTHORIZED","message":"Invalid or expired token"}}`, http.StatusUnauthorized)
				return
			}

			uid, err := uuid.Parse(claims.UserID)
			if err != nil {
				http.Error(w, `{"success":false,"error":{"code":"UNAUTHORIZED","message":"Invalid user identity in token"}}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDContextKey, uid)
			ctx = context.WithValue(ctx, UserRoleContextKey, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	val, ok := ctx.Value(UserIDContextKey).(uuid.UUID)
	return val, ok
}
