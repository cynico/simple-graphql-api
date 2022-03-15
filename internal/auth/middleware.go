package auth

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/cynico/hackernews/internal/users"
	"github.com/cynico/hackernews/pkg/jwt"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			// Allow unauthenticated users in
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			tokenStr := header
			username, err := jwt.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				log.Printf("%v", err)
				return
			}

			user := users.User{Username: username}
			id, err := users.GetUserIdByUsername(user.Username)
			if err != nil {
				next.ServeHTTP(w, r)
			}
			user.ID = strconv.Itoa(id)

			ctx := context.WithValue(r.Context(), userCtxKey, &user)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *users.User {
	raw, _ := ctx.Value(userCtxKey).(*users.User)
	return raw
}
