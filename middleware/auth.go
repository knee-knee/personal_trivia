package middleware

import (
	"log"
	"net/http"

	"github.com/personal_trivia/repo"
)

const AuthHeader = "user-id"

func (m *Middleware) AuthMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println("middleware: Trying to find user")
		// Right now this is just the user ID till I have a better auth system.
		token := req.Header.Get(AuthHeader)
		if token == "" {
			http.Error(w, "Did not provide auth header", http.StatusBadRequest)
			return
		}

		user, err := m.Repo.GetUser(token)
		if err != nil {
			http.Error(w, "Could not retrieve user from dynamo", http.StatusInternalServerError)
			return
		}

		if (user != repo.User{}) {
			log.Printf("middleware: User %s successfully authed. \n", token)
			next.ServeHTTP(w, req)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}
