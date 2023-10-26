package model

import "quick_web_golang/provider"

type Repo struct {
	UserRepo *UserRepo
}

var Repos = &Repo{}

func NewRepo() *Repo {
	return &Repo{
		UserRepo: NewUserRepo(provider.Database.DB),
	}
}
