package middlewares

import (
	"TaskList/internal/lib/jwt"
	"context"
	"github.com/go-chi/render"
	"net/http"
	"strings"
)

const (
	prefix = "Bearer "
)

func AuthMW(secretKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//Получить токен
			bearer := r.Header.Get("Authorization")
			if bearer == "" || !strings.HasPrefix(bearer, prefix) {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, "Отсутсвует заголовок авторизации")
				return
			}
			token := strings.TrimPrefix(bearer, prefix)
			//Проверить валидность токена
			claims, err := jwt.ValidateToken(token, []byte(secretKey))
			if err != nil {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, err.Error())
				return
			}
			//Добавляем claims в контекст
			r = r.WithContext(context.WithValue(r.Context(), "claims", claims))
			//Пропускаем запрос далее
			next.ServeHTTP(w, r)
		})
	}
}
