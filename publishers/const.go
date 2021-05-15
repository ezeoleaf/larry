package publishers

type IPublish interface {
	PublishContent(s string) bool
}

// Twitter is the value of the valid publisher
const Twitter = "twitter"

var ValidPublishers = []string{
	Twitter,
}

func GetValidPublishersToString() string {
	p := ""

	for _, v := range ValidPublishers {
		if p != "" {
			p += ", "
		}

		p += v
	}

	return "The valid publishers are: " + p
}
