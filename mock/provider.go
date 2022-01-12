package mock

import "github.com/ezeoleaf/larry/domain"

// ProviderMock is a mock Provider
type ProviderMock struct {
	GetContentToPublishFn func() (*domain.Content, error)
}

// GetContentToPublish calls the GetContentToPublishFn
func (p ProviderMock) GetContentToPublish() (*domain.Content, error) {
	if p.GetContentToPublishFn == nil {
		return nil, nil
	}

	return p.GetContentToPublishFn()
}
