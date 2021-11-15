package mock

// PublisherMock is a mock Publisher
type PublisherMock struct {
	PublishContentFn func(content string) (bool, error)
}

// PublishContent calls the PublishContentFn
func (p PublisherMock) PublishContent(content string) (bool, error) {
	if p.PublishContentFn == nil {
		return false, nil
	}

	return p.PublishContentFn(content)
}
