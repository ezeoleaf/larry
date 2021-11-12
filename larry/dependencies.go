package larry

type Provider interface {
	GetContentToPublish() (string, error)
}

type Publisher interface {
	PublishContent(content string) (bool, error)
}
