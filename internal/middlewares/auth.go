package middlewares

import (
	"TaskList/internal/lib/jwt"
	"context"
	"github.com/go-chi/render"
	"net/http"
	"strings"
	"time"
)

const (
	prefix = "Bearer "
)

type Key string

const KeyClaims Key = "claims"

func AuthJWT(secretKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//get token
			bearer := r.Header.Get("Authorization")
			if bearer == "" || !strings.HasPrefix(bearer, prefix) {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, "Отсутсвует заголовок авторизации")
				return
			}
			token := strings.TrimPrefix(bearer, prefix)
			//check valid token
			claims, err := jwt.ValidateToken(token, []byte(secretKey))
			if err != nil {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, err.Error())
				return
			}

			// check exp time token
			t, err := claims.GetExpirationTime()
			if !t.After(time.Now()) {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, err.Error())
				return
			}

			//add claims in context
			r = r.WithContext(context.WithValue(r.Context(), KeyClaims, claims))

			next.ServeHTTP(w, r)
		})
	}
}
