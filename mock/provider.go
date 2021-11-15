package mock

// ProviderMock is a mock Provider
type ProviderMock struct {
	GetContentToPublishFn func() (string, error)
}

// GetContentToPublish calls the GetContentToPublishFn
func (p ProviderMock) GetContentToPublish() (string, error) {
	if p.GetContentToPublishFn == nil {
		return "", nil
	}

	return p.GetContentToPublishFn()
}
