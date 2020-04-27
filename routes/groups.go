package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/personal_trivia/repo"
)

const AuthHeader = "user-id"

type Group struct {
	ID      *string `json:"ID,omitempty"`
	Name    *string `json:"name,omitempty"`
	Creator *string `json:"creator,omitempty"`
}

func (g Group) ToDynamo() repo.Group {
	return repo.Group{
		ID:      aws.StringValue(g.ID),
		Name:    aws.StringValue(g.Name),
		Creator: aws.StringValue(g.Creator),
	}
}

func (r *Routes) CreateGroup(w http.ResponseWriter, req *http.Request) {
	log.Println("routes: Strarting to create group")
	var in Group
	if err := json.NewDecoder(req.Body).Decode(&in); err != nil {
		http.Error(w, "could not unmarshal input", 500)
		return
	}

	defer req.Body.Close()

	// We know this will not be empty because this is an authenticated route.
	creator := req.Header.Get(AuthHeader)
	in.Creator = &creator

	resp, err := r.Repo.CreateGroup(in.ToDynamo())
	if err != nil {
		http.Error(w, "could not create user", 500)
		return
	}

	log.Printf("routes: Successfully created group with ID %s \n", resp)

	w.Write([]byte(resp))
}
