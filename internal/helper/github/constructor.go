package github

import (
	"github.com/google/go-github/github"
	"github.com/nori-io/nori/internal/domain/helper"
)

type Helper struct {
	client *github.Client
}

type Params struct {
	Client *github.Client
}

func New(params Params) helper.GithubHelperInterface {
	return &Helper{
		client: params.Client,
	}
}
