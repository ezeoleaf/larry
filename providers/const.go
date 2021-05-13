package providers

type IContent interface {
	GetContentToPublish() string
}

// Github is the value of the valid provider
const Github = "github"

var ValidProviders = []string{
	Github,
}

func GetValidProvidersToString() string {
	p := ""

	for _, v := range ValidProviders {
		if p != "" {
			p += ", "
		}

		p += v
	}

	return "The valid providers are: " + p
}
