package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/personal_trivia/repo"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
	Group    string `json:"group"`
}

func (u User) ToDynamo() repo.User {
	return repo.User{
		Email:    u.Email,
		Password: u.Password,
		Username: u.Username,
		Group:    u.Group,
	}
}

// TODO: we need to check first that the username and email are unique.
func (r *Routes) CreateUser(w http.ResponseWriter, req *http.Request) {
	log.Println("routes: Starting to create user")
	var in User
	if err := json.NewDecoder(req.Body).Decode(&in); err != nil {
		http.Error(w, "could not unmarshal input", 500)
		return
	}

	defer req.Body.Close()

	resp, err := r.Repo.CreateUser(in.ToDynamo())
	if err != nil {
		http.Error(w, "could not create user", 500)
		return
	}

	w.Write([]byte(resp))
}
