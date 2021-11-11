package larry

import (
	"fmt"
	"log"
)

type publisher interface {
	PublishContent(content string) error
}

type provider interface {
	GetContentToPublish() (string, error)
}

type config interface {
}

// Service braze service.
type Service struct {
	Publishers map[string]publisher
	Provider   provider
	Config     config
	Logger     log.Logger
}

func (s Service) Run() error {

	content, err := s.Provider.GetContentToPublish()
	if err != nil {
		return err
	}

	fmt.Println(content)

	return nil
}
