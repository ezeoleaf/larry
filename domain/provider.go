package domain

type Client interface {
	GetContentToPublish() (string, error)
}
