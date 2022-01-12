package larry

import "github.com/ezeoleaf/larry/domain"

// Provider represents the interface for the different providers
type Provider interface {
	GetContentToPublish() (*domain.Content, error)
}

// Publisher represents the interface for the different publishers
type Publisher interface {
	PublishContent(content *domain.Content) (bool, error)
}
