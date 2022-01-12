package mock

import "github.com/ezeoleaf/larry/domain"

// PublisherMock is a mock Publisher
type PublisherMock struct {
	PublishContentFn func(content *domain.Content) (bool, error)
}

// PublishContent calls the PublishContentFn
func (p PublisherMock) PublishContent(content *domain.Content) (bool, error) {
	if p.PublishContentFn == nil {
		return false, nil
	}

	return p.PublishContentFn(content)
}
