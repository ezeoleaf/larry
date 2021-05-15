package providers

import "testing"

func TestConsts(t *testing.T) {
	if Github != "github" {
		t.Errorf("The github provider is wrong. Expected github got %s", Github)
	}
}

func TestGetValidProvidersToString(t *testing.T) {
	ValidProviders = []string{
		"a",
	}

	r := GetValidProvidersToString()
	expected := "The valid providers are: a"

	if r != expected {
		t.Errorf("The provides are wrong. Expected: %s got %s", expected, r)
	}

	ValidProviders = []string{
		"a",
		"b",
	}

	r = GetValidProvidersToString()
	expected = "The valid providers are: a, b"

	if r != expected {
		t.Errorf("The provides are wrong. Expected: %s got %s", expected, r)
	}
}
