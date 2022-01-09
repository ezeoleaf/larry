package github

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"

	"github.com/ezeoleaf/larry/config"
	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

type client interface {
	GetReadme(ctx context.Context, owner, repo string, opts *github.RepositoryContentGetOptions) (*github.RepositoryContent, *github.Response, error)
	GetContents(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error)
}

type repositoryData struct {
	Owner string
	Name  string
}

// Publisher represents the publisher client
type Publisher struct {
	GithubClient   client
	Config         config.Config
	RepositoryData repositoryData
}

const repoFileName = "README.md"

// NewPublisher returns a new publisher client
func NewPublisher(apiKey string, cfg config.Config, repoOwner, repoName string) Publisher {
	log.Print("New Github Publisher")
	p := Publisher{Config: cfg, RepositoryData: repositoryData{Owner: repoOwner, Name: repoName}}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apiKey},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	p.GithubClient = github.NewClient(tc).Repositories

	return p
}

func (p Publisher) getReadmeContent() (string, error) {
	ctx := context.Background()
	fc, _, _, err := p.GithubClient.GetContents(ctx, p.RepositoryData.Owner, p.RepositoryData.Name, repoFileName, nil)

	if err != nil {
		fmt.Println(fmt.Errorf("could not fetch repository, got %v", err))
		return "", err
	}

	if fc == nil || fc.Content == nil {
		fmt.Println("content of README is empty")
		return "", errors.New("no content")
	}

	rawDecodedText, err := base64.StdEncoding.DecodeString(*fc.Content)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decoded text: %s\n", rawDecodedText)

	return "s", nil
}

// PublishContent receives a content to publish and try to publish in README file
func (p Publisher) PublishContent(content string) (bool, error) {
	_, _ = p.getReadmeContent()
	if p.Config.SafeMode {
		log.Print("Running in Safe Mode")
		log.Print(content)
		fmt.Println(content)
		return true, nil
	}

	return true, nil

}
