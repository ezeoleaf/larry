package larry

import (
	"log"
)

// Service represents the application struct
type Service struct {
	Publishers map[string]Publisher
	Provider   Provider
	Logger     *log.Logger
}

// Run executes the application
func (s Service) Run() error {
	content, err := s.Provider.GetContentToPublish()
	if err != nil {
		s.Logger.Println("error on getting content for publishing", err)
		return err
	}

	if content != nil {
		for _, pub := range s.Publishers {
			_, err := pub.PublishContent(content)
			if err != nil {
				s.Logger.Println("error on publishing", err)
				return err
			}
		}
	}

	return nil
}
