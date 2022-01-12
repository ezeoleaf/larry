package github

import (
	"context"
	"errors"
	"testing"

	"github.com/ezeoleaf/larry/config"
	"github.com/ezeoleaf/larry/domain"
	"github.com/ezeoleaf/larry/mock"
	"github.com/google/go-github/v39/github"
)

func TestNewPublisher(t *testing.T) {
	c := config.Config{SafeMode: true}
	apiKey := "s"

	p := NewPublisher(apiKey, c, "o", "n", "fn")

	if p.GithubClient == nil {
		t.Error("expected new publisher, got nil")
	}
}

func TestGetReadmeContent(t *testing.T) {
	fileName := "fileName"
	fileContentStr := "content"
	for _, tc := range []struct {
		Name        string
		repoClient  mock.RepoClientMock
		returnValue *github.RepositoryContent
		shouldFail  bool
	}{
		{
			Name: "Test error file content",
			repoClient: mock.RepoClientMock{
				GetContentsFn: func(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error) {
					return nil, nil, nil, errors.New("could not retrieve content")
				},
			},
			returnValue: nil,
			shouldFail:  true,
		},
		{
			Name: "Test no error but no file content",
			repoClient: mock.RepoClientMock{
				GetContentsFn: func(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error) {
					fc := github.RepositoryContent{Name: &fileName}
					return &fc, nil, nil, nil
				},
			},
			returnValue: nil,
			shouldFail:  true,
		},
		{
			Name: "Test returns file content",
			repoClient: mock.RepoClientMock{
				GetContentsFn: func(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error) {
					fc := github.RepositoryContent{Name: &fileName, Content: &fileContentStr}
					return &fc, nil, nil, nil
				},
			},
			returnValue: &github.RepositoryContent{Name: &fileName, Content: &fileContentStr},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			p := Publisher{GithubClient: tc.repoClient}
			resp, err := p.getReadmeContent(context.Background())

			if tc.shouldFail && err == nil {
				t.Error("expected error but got nil instead")
			} else if !tc.shouldFail && err != nil {
				t.Errorf("expected no error but got %v instead", err)
			}

			if resp == nil && resp != tc.returnValue {
				t.Errorf("expected not content but got %v instead", tc.returnValue)
			} else if resp != nil {
				if *resp.Name != *tc.returnValue.Name {
					t.Errorf("expected file name %s but got %s instead", *resp.Name, *tc.returnValue.Name)
				}

				if *resp.Content != *tc.returnValue.Content {
					t.Errorf("expected content %s but got %s instead", *resp.Content, *tc.returnValue.Content)
				}
			}
		})
	}
}

func TestDecodeBase64(t *testing.T) {
	encodedText := "c29tZXRleHQ="

	for _, tc := range []struct {
		Name          string
		expectedValue string
		shouldFail    bool
		text          string
	}{
		{
			Name:          "Test could not decode",
			text:          "some text",
			expectedValue: "",
			shouldFail:    true,
		},
		{
			Name:          "Test should decode",
			expectedValue: "sometext",
			text:          encodedText,
			shouldFail:    false,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			resp, err := decodeBase64(tc.text)

			if tc.shouldFail && err == nil {
				t.Error("expected error but got nil instead")
			} else if !tc.shouldFail && err != nil {
				t.Errorf("expected no error but got %v instead", err)
			}

			if resp != tc.expectedValue {
				t.Errorf("expected %s but got %s instead", tc.expectedValue, resp)
			}
		})
	}
}

func TestPublishContent(t *testing.T) {
	ti, sub, url := "sometext", "sub", "url"
	// fileName := "fileName"
	// fileContentStr := "content"
	encodedText := "c29tZXRleHQ="
	for _, tc := range []struct {
		Name          string
		repoClient    mock.RepoClientMock
		expectedValue bool
		shouldFail    bool
		content       *domain.Content
		cfg           config.Config
	}{
		{
			Name:          "Test no content to publish error",
			expectedValue: false,
			shouldFail:    true,
			cfg:           config.Config{SafeMode: true},
		},
		{
			Name:          "Test running in safe mode without some values",
			expectedValue: false,
			shouldFail:    true,
			cfg:           config.Config{SafeMode: true},
			content:       &domain.Content{Title: &ti},
		},
		{
			Name:          "Test running in safe mode without subtitle",
			expectedValue: true,
			shouldFail:    false,
			cfg:           config.Config{SafeMode: true},
			content:       &domain.Content{Title: &ti, URL: &url},
		},
		{
			Name:          "Test running in safe mode with subtitle",
			expectedValue: true,
			shouldFail:    false,
			cfg:           config.Config{SafeMode: true},
			content:       &domain.Content{Title: &ti, Subtitle: &sub, URL: &url},
		},
		{
			Name:          "Test fail reading readme content",
			expectedValue: false,
			shouldFail:    true,
			cfg:           config.Config{},
			content:       &domain.Content{Title: &ti, Subtitle: &sub, URL: &url},
			repoClient: mock.RepoClientMock{
				GetContentsFn: func(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error) {
					return nil, nil, nil, errors.New("could not retrieve content")
				},
			},
		},
		{
			Name:          "Test fail decoding readme",
			expectedValue: false,
			shouldFail:    true,
			cfg:           config.Config{},
			content:       &domain.Content{Title: &ti, Subtitle: &sub, URL: &url},
			repoClient: mock.RepoClientMock{
				GetContentsFn: func(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error) {
					someText := "some text"
					fc := github.RepositoryContent{Content: &someText}
					return &fc, nil, nil, nil
				},
			},
		},
		{
			Name:          "Test repo already added",
			expectedValue: false,
			shouldFail:    true,
			cfg:           config.Config{},
			content:       &domain.Content{Title: &ti, Subtitle: &sub, URL: &url},
			repoClient: mock.RepoClientMock{
				GetContentsFn: func(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error) {
					fc := github.RepositoryContent{Content: &encodedText}
					return &fc, nil, nil, nil
				},
			},
		},
		{
			Name:          "Test error on publish repository",
			expectedValue: false,
			shouldFail:    true,
			cfg:           config.Config{},
			content:       &domain.Content{Title: &ti, Subtitle: &sub, URL: &url},
			repoClient: mock.RepoClientMock{
				GetContentsFn: func(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error) {
					otherText := "sometext"
					fc := github.RepositoryContent{Content: &otherText}
					return &fc, nil, nil, nil
				},
				UpdateFileFn: func(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error) {
					return nil, nil, errors.New("could not publish")
				},
			},
		},
		{
			Name:          "Test should publish repository",
			expectedValue: true,
			shouldFail:    false,
			cfg:           config.Config{},
			content:       &domain.Content{Title: &ti, Subtitle: &sub, URL: &url},
			repoClient: mock.RepoClientMock{
				GetContentsFn: func(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error) {
					otherText := "sometext"
					fc := github.RepositoryContent{Content: &otherText}
					return &fc, nil, nil, nil
				},
				UpdateFileFn: func(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error) {
					return nil, nil, nil
				},
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			p := Publisher{GithubClient: tc.repoClient, Config: tc.cfg}
			resp, err := p.PublishContent(tc.content)

			if tc.shouldFail && err == nil {
				t.Error("expected error but got nil instead")
			} else if !tc.shouldFail && err != nil {
				t.Errorf("expected no error but got %v instead", err)
			}

			if resp != tc.expectedValue {
				t.Errorf("expected %v as response but got %v instead", tc.expectedValue, resp)
			}
		})
	}
}
