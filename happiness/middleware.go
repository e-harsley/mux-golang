package happiness

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type MiddlewareFunc func(http.Handler) http.Handler

type Middleware struct {
	PublicRoute bool
}

// AuthDeps is a middleware function for authentication.
func (m Middleware) AuthDeps(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.PublicRoute {
			// If the route is public, don't perform authentication.
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid token format", http.StatusBadRequest)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		auth := AuthToken{}
		authToken, err := auth.ParseAuthToken(token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("did i get here", authToken)
		ctx := context.WithValue(r.Context(), "user_context", ExtraParameters{UserID: authToken.UserID})
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})

}
