package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/personal_trivia/repo"
)

type Question struct {
	Question      string `json:"question"`
	AnwserA       string `json:"anwserA"`
	AnwserB       string `json:"anwserB"`
	AnwserC       string `json:"anwserC"`
	AnwserD       string `json:"anwserD"`
	CorrectAnwser string `json:"correct_anwser,omitempty"`
}

func (q Question) ToDynamo() repo.Question {
	return repo.Question{
		Question: q.Question,
		AnwserA:  q.AnwserA,
		AnwserB:  q.AnwserB,
		AnwserC:  q.AnwserC,
		AnwserD:  q.AnwserD,
	}
}

func questionFromDynamo(resp repo.Question) Question {
	return Question{
		Question:      resp.Question,
		AnwserA:       resp.AnwserA,
		AnwserB:       resp.AnwserB,
		AnwserC:       resp.AnwserC,
		AnwserD:       resp.AnwserD,
		CorrectAnwser: resp.CorrectAnwser,
	}
}

func (r *Routes) Question(w http.ResponseWriter, req *http.Request) {
	log.Println("routes: Starting to fetch question")
	params := mux.Vars(req)
	questionID := params["id"]
	resp, err := r.Repo.GetQuestion(questionID)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	if (resp == repo.Question{}) {
		http.Error(w, "could not find question", 404)
	}

	question := questionFromDynamo(resp)
	body, err := json.Marshal(question)
	if err != nil {
		http.Error(w, "could not marshal response from dynamo", 500)
	}

	w.Write(body)
}

func (r *Routes) CreateQuestion(w http.ResponseWriter, req *http.Request) {
	log.Println("routes: Starting to create a question")
	var in Question
	if err := json.NewDecoder(req.Body).Decode(&in); err != nil {
		http.Error(w, "could not unmarshal input", 500)
	}

	defer req.Body.Close()

	resp, err := r.Repo.CreateQuestion(in.ToDynamo())
	if err != nil {
		http.Error(w, "could not create the question", 500)
	}

	w.Write([]byte(resp))
}
