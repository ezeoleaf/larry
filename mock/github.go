package mock

import (
	"context"

	"github.com/google/go-github/v39/github"
)

type SearchClientMock struct {
	RepositoriesFn func(ctx context.Context, query string, opt *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error)
}

type UserClientMock struct {
	GetFn func(ctx context.Context, user string) (*github.User, *github.Response, error)
}

func (scm SearchClientMock) Repositories(ctx context.Context, query string, opt *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error) {
	if scm.RepositoriesFn == nil {
		return nil, nil, nil
	}

	return scm.RepositoriesFn(ctx, query, opt)
}

func (ucm UserClientMock) Get(ctx context.Context, user string) (*github.User, *github.Response, error) {
	if ucm.GetFn == nil {
		return nil, nil, nil
	}

	return ucm.GetFn(ctx, user)
}
