package middleware

import "github.com/personal_trivia/repo"

type Middleware struct {
	Repo *repo.Repo
}

func New(r *repo.Repo) *Middleware {
	return &Middleware{
		Repo: r,
	}
}
