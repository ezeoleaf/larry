package larry

import (
	"log"
)

// Service represents the application struct
type Service struct {
	Publishers map[string]Publisher
	Provider   Provider
	Logger     log.Logger
}

// Run executes the application
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
