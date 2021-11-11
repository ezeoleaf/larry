package larry

import (
	"log"
)

type config interface {
}

// Service braze service.
type Service struct {
	Publishers map[string]Publisher
	Provider   Provider
	Config     config
	Logger     log.Logger
}

func (s Service) Run() error {

	content, err := s.Provider.GetContentToPublish()
	if err != nil {
		return err
	}

	for _, pub := range s.Publishers {
		pub.PublishContent(content)
	}

	return nil
}
