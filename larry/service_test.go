package larry

import (
	"errors"
	"testing"

	"github.com/ezeoleaf/larry/mock"
)

func TestRun(t *testing.T) {
	for _, tc := range []struct {
		Name           string
		mockProvider   mock.ProviderMock
		mockPublishers map[string]Publisher
		returnValue    error
	}{
		{
			Name: "Test error on get content",
			mockProvider: mock.ProviderMock{
				GetContentToPublishFn: func() (string, error) {
					return "", errors.New("some error")
				},
			},
			mockPublishers: map[string]Publisher{
				"mock": mock.PublisherMock{},
			},
			returnValue: errors.New("some error"),
		},
		{
			Name: "Test get content and publish",
			mockProvider: mock.ProviderMock{
				GetContentToPublishFn: func() (string, error) {
					return "content", nil
				},
			},
			mockPublishers: map[string]Publisher{
				"mock": mock.PublisherMock{
					PublishContentFn: func(content string) (bool, error) {
						return true, nil
					},
				},
			},
			returnValue: nil,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			s := Service{
				Provider:   tc.mockProvider,
				Publishers: tc.mockPublishers,
			}

			resp := s.Run()

			if tc.returnValue != nil {
				if resp == nil {
					t.Errorf("expected resp but got %v instead", resp)
					return
				}

				if tc.returnValue.Error() != resp.Error() {
					t.Errorf("expected %v as value, got %v instead", tc.returnValue.Error(), resp.Error())
					return
				}
			} else if tc.returnValue == nil && resp != nil {
				t.Errorf("expected no resp but got %v instead", resp)
			}

		})
	}

}
