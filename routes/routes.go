package routes

import (
	"github.com/personal_trivia/repo"
)

type Routes struct {
	Repo *repo.Repo
}

func New(r *repo.Repo) *Routes {
	return &Routes{
		Repo: r,
	}
}
