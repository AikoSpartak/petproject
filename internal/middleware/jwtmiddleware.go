package middleware

import (
	"arch/internal/services"
	"context"
	"net/http"
	"strings"
)

const userIDKey = "user_id"
const userRole = "role"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Authorization")
		tokenStr := strings.TrimPrefix(header, "Bearer ")

		if tokenStr == header || tokenStr == "" {
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}

		claims, err := services.ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(
			r.Context(),
			userIDKey,
			claims.UserID,
		)

		ctx = context.WithValue(
			r.Context(),
			userRole,
			claims.Role,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}

func UserRoleFromContext(ctx context.Context) (string, bool) {

	role, ok := ctx.Value(userRole).(string)
	return role, ok
}
