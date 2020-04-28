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

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u User) ToDynamo() repo.User {
	return repo.User{
		Email:    u.Email,
		Password: u.Password,
		Username: u.Username,
		Group:    u.Group,
	}
}

// Rules for creating a user:
// 	1. Email is unique
// 	2. Username is unique
// 	3. Group is a valid group ID.
// TODO: we need to check first that the username and email are unique.
func (r *Routes) CreateUser(w http.ResponseWriter, req *http.Request) {
	log.Println("routes: Starting to create user")
	var in User
	if err := json.NewDecoder(req.Body).Decode(&in); err != nil {
		http.Error(w, "could not unmarshal input", 500)
		return
	}

	defer req.Body.Close()

	// Check if email is unique.
	user, err := r.Repo.GetUserScan(repo.UserScanInput{
		Key:   "email",
		Value: in.Email,
	})
	if err != nil {
		log.Printf("routes: could not scan for user based off email %v", err)
		http.Error(w, "internall server error", 500)
		return
	}
	if (user != repo.User{}) {
		http.Error(w, "user with email already exists", 400)
		return
	}

	// Check if username is unique.
	user, err = r.Repo.GetUserScan(repo.UserScanInput{
		Key:   "username",
		Value: in.Username,
	})
	if err != nil {
		log.Printf("routes: could not scan for user based off username %v \n", err)
		http.Error(w, "internall server error", 500)
		return
	}
	if (user != repo.User{}) {
		http.Error(w, "user with username already exists", 400)
		return
	}

	// Check if group is valid.
	group, err := r.Repo.GetGroup(in.Group)
	if err != nil {
		log.Printf("routes: could not retrieve group from dynamo %v \n", err)
		http.Error(w, "internall server error", 500)
		return
	}
	if (group == repo.Group{}) {
		http.Error(w, "Trying to create a user in an invalid group", http.StatusBadRequest)
		return
	}

	resp, err := r.Repo.CreateUser(in.ToDynamo())
	if err != nil {
		http.Error(w, "could not create user", 500)
		return
	}

	log.Printf("routes: Successfully created user with ID %s \n", resp)

	w.Write([]byte(resp))
}

func (r *Routes) Login(w http.ResponseWriter, req *http.Request) {
	var in Login
	if err := json.NewDecoder(req.Body).Decode(&in); err != nil {
		http.Error(w, "could not unmarshal login input", 500)
		return
	}

	defer req.Body.Close()

	if in.Email == "" || in.Password == "" {
		http.Error(w, "invalid email or password", http.StatusBadRequest)
		return
	}

	log.Printf("routes: Attempting to login user %s \n", in.Email)

	resp, err := r.Repo.GetUserScan(repo.UserScanInput{
		Key:   "email",
		Value: in.Email,
	})
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusBadRequest)
		return
	}

	if resp.Password != in.Password {
		http.Error(w, "invalid email or password", http.StatusBadRequest)
		return
	}

	log.Printf("routes: Successfully logged in user %s \n", in.Email)

	w.Write([]byte(resp.ID))
}
