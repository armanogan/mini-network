package transport

import (
	"mini-network/internal/v1/models"
	"mini-network/internal/v1/token"
	"net/http"
)

type Midleware struct {
	handler func(http.ResponseWriter, *http.Request)
}

func Middleware(next func(http.ResponseWriter, *http.Request, models.User), jwt token.TokenService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("access_token")
		var user models.User
		if err == nil && cookie.Valid() == nil {
			user, err = jwt.ParseToken(cookie.Value)
		}
		if err != nil {
			http.Redirect(w, req, "/login/", http.StatusMovedPermanently)
			return
		}
		next(w, req, user)
	})
}
