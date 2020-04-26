package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/personal_trivia/repo"
)

type Answer struct {
	Answer string `json:"answer"`
}

type Question struct {
	Question      string `json:"question"`
	AnswerA       string `json:"answerA"`
	AnswerB       string `json:"answerB"`
	AnswerC       string `json:"answerC"`
	AnswerD       string `json:"answerD"`
	CorrectAnswer string `json:"correct_answer,omitempty"`
}

func (q Question) ToDynamo() repo.Question {
	return repo.Question{
		Question:      q.Question,
		AnswerA:       q.AnswerA,
		AnswerB:       q.AnswerB,
		AnswerC:       q.AnswerC,
		AnswerD:       q.AnswerD,
		CorrectAnswer: q.CorrectAnswer,
	}
}

func questionFromDynamo(resp repo.Question) Question {
	return Question{
		Question:      resp.Question,
		AnswerA:       resp.AnswerA,
		AnswerB:       resp.AnswerB,
		AnswerC:       resp.AnswerC,
		AnswerD:       resp.AnswerD,
		CorrectAnswer: resp.CorrectAnswer,
	}
}

func (r *Routes) Question(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	questionID, ok := params["id"]
	if !ok {
		http.Error(w, "did not provide question ID", 400)
		return
	}
	log.Printf("routes: Starting to fetch question %s \n", questionID)
	resp, err := r.Repo.GetQuestion(questionID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if (resp == repo.Question{}) {
		http.Error(w, "could not find question", 404)
		return
	}

	question := questionFromDynamo(resp)
	body, err := json.Marshal(question)
	if err != nil {
		http.Error(w, "could not marshal response from dynamo", 500)
		return
	}

	w.Write(body)
}

func (r *Routes) CheckAnswer(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	questionID, ok := params["id"]
	if !ok {
		http.Error(w, "did not provide question ID", 400)
		return
	}
	log.Printf("routes: Starting to fetch question %s \n", questionID)
	resp, err := r.Repo.GetQuestion(questionID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if (resp == repo.Question{}) {
		http.Error(w, "could not find question", 404)
		return
	}
	log.Println("routes: checking if answer is correct")

	var in Answer
	if err := json.NewDecoder(req.Body).Decode(&in); err != nil {
		http.Error(w, "could not unmarshal input", 500)
		return
	}

	defer req.Body.Close()

	if resp.CorrectAnswer == in.Answer {
		w.Write([]byte("true"))
	} else {
		w.Write([]byte("false"))
	}
}

func (r *Routes) CreateQuestion(w http.ResponseWriter, req *http.Request) {
	log.Println("routes: Starting to create a question")
	var in Question
	if err := json.NewDecoder(req.Body).Decode(&in); err != nil {
		http.Error(w, "could not unmarshal input", 500)
		return
	}

	defer req.Body.Close()

	resp, err := r.Repo.CreateQuestion(in.ToDynamo())
	if err != nil {
		http.Error(w, "could not create the question", 500)
		return
	}

	w.Write([]byte(resp))
}
