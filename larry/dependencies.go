package larry

// Provider represents the interface for the different providers
type Provider interface {
	GetContentToPublish() (string, error)
}

// Publisher represents the interface for the different publishers
type Publisher interface {
	PublishContent(content string) (bool, error)
}
