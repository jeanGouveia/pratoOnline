package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/jeanGouveia/pratoOnline/backend/internal/service"
)

type contextKey string

const ContextKeyUserID contextKey = "user_id"
const ContextKeyClaims contextKey = "claims"

type AuthMiddleware struct {
	authService *service.AuthService
}

func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (m *AuthMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		// Estratégia 1: Cookie HttpOnly (produção)
		if cookie, err := r.Cookie("auth_token"); err == nil {
			token = cookie.Value
		}

		// Estratégia 2: Authorization header (dev / Postman)
		if token == "" {
			if h := r.Header.Get("Authorization"); strings.HasPrefix(h, "Bearer ") {
				token = strings.TrimPrefix(h, "Bearer ")
			}
		}

		if token == "" {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		claims, err := m.authService.ValidateToken(token)
		if err != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		// Injeta UserID e claims completos no context
		ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextKeyClaims, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext extrai o UserID injetado pelo middleware.
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	id, ok := ctx.Value(ContextKeyUserID).(uint)
	return id, ok
}

// GetClaimsFromContext extrai os claims completos (UserID, Email, Name).
func GetClaimsFromContext(ctx context.Context) (*service.JWTClaims, bool) {
	claims, ok := ctx.Value(ContextKeyClaims).(*service.JWTClaims)
	return claims, ok
}