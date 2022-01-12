package mock

import (
	"context"

	"github.com/google/go-github/v39/github"
)

// SearchClientMock is the mock for the search client from Github provider
type SearchClientMock struct {
	RepositoriesFn func(ctx context.Context, query string, opt *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error)
}

// UserClientMock is the mock for the user client from Github provider
type UserClientMock struct {
	GetFn func(ctx context.Context, user string) (*github.User, *github.Response, error)
}

type RepoClientMock struct {
	GetContentsFn func(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error)
	UpdateFileFn  func(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error)
}

// Repositories calls RepositoriesFn
func (scm SearchClientMock) Repositories(ctx context.Context, query string, opt *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error) {
	if scm.RepositoriesFn == nil {
		return nil, nil, nil
	}

	return scm.RepositoriesFn(ctx, query, opt)
}

// Get calls GetFn
func (ucm UserClientMock) Get(ctx context.Context, user string) (*github.User, *github.Response, error) {
	if ucm.GetFn == nil {
		return nil, nil, nil
	}

	return ucm.GetFn(ctx, user)
}

func (rcm RepoClientMock) GetContents(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error) {
	if rcm.GetContentsFn == nil {
		return nil, nil, nil, nil
	}
	return rcm.GetContentsFn(ctx, owner, repo, path, opts)
}
func (rcm RepoClientMock) UpdateFile(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error) {
	if rcm.UpdateFileFn == nil {
		return nil, nil, nil
	}
	return rcm.UpdateFileFn(ctx, owner, repo, path, opts)
}
