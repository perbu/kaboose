package localgit

import "github.com/go-git/go-git/v5"

// LocalRepo is a wrapper around the go-git client
type LocalRepo struct {
	path string
	git  *git.Repository
}

// NewClient creates a new go-git client
func NewLocalRepo(path string) (*LocalRepo, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}
	c := &LocalRepo{
		path: path,
		git:  repo,
	}
	return c, nil
}
