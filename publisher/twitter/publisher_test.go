package twitter

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/ezeoleaf/larry/config"
// 	"github.com/ezeoleaf/larry/domain"
// )

// func TestNewPublisher(t *testing.T) {
// 	c := config.Config{SafeMode: true}
// 	ak := AccessKeys{}

// 	p := NewPublisher(ak, c)

// 	if p.Client == nil {
// 		t.Error("expected new publisher, got nil")
// 	}
// }

// func TestPublishContentInSafeMode(t *testing.T) {
// 	c := config.Config{SafeMode: true}
// 	ak := AccessKeys{}

// 	p := NewPublisher(ak, c)

// 	ti, s, u := "ti", "s", "u"

// 	cont := domain.Content{Title: &ti, Subtitle: &s, URL: &u}

// 	r, err := p.PublishContent(&cont)

// 	if !r {
// 		t.Error("expected content published in Safe Mode. No content published")
// 	}

// 	if err != nil {
// 		t.Errorf("expected no error got %v", err)
// 	}
// }

// func TestCheckTweetDataInSafeMode(t *testing.T) {
// 	c := config.Config{SafeMode: true}
// 	ak := AccessKeys{}

// 	p := NewPublisher(ak, c)

// 	ti, u := "Lorem Ipsum", "https://loremipsum.io/generator/?n=3&t=s"
// 	extraData := []string{"50k", "Author: @unknown"}

// 	for _, tc := range []struct {
// 		Name           string
// 		Subtitle       string
// 		ExpectedResult string
// 	}{
// 		{
// 			Name:           "Test should return same content",
// 			Subtitle:       "t",
// 			ExpectedResult: "Lorem Ipsum: t\n50k\nAuthor: @unknown\nhttps://loremipsum.io/generator/?n=3&t=s",
// 		},
// 		{
// 			Name:           "Test should truncate subtitle",
// 			Subtitle:       "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Vitae sapien pellentesque habitant morbi tristique senectus et netus et. Nunc sed velit dignissim sodales.",
// 			ExpectedResult: "Lorem Ipsum: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Vitae sapien pellentesque habitant morbi tristique senectus et netus et. Nunc ...\n50k\nAuthor: @unknown\nhttps://loremipsum.io/generator/?n=3&t=s",
// 		},
// 	} {
// 		t.Run(tc.Name, func(t *testing.T) {
// 			cont := &domain.Content{Title: &ti, Subtitle: &tc.Subtitle, URL: &u, ExtraData: extraData}

// 			resp := p.prepareTweet(cont)
// 			fmt.Println(resp)
// 			if resp != tc.ExpectedResult {
// 				t.Errorf("resp should be %v, got %v", tc.ExpectedResult, resp)
// 			}

// 			if len(resp) > TweetLength {
// 				t.Errorf("tweet length is %v, which is greater than %v", len(resp), TweetLength)
// 			}
// 		})
// 	}

// }
