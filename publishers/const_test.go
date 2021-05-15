package publishers

import "testing"

func TestConsts(t *testing.T) {
	if Twitter != "twitter" {
		t.Errorf("The twitter provider is wrong. Expected twitter got %s", Twitter)
	}
}

func TestGetValidPublishersToString(t *testing.T) {
	ValidPublishers = []string{
		"a",
	}

	r := GetValidPublishersToString()
	expected := "The valid publishers are: a"

	if r != expected {
		t.Errorf("The provides are wrong. Expected: %s got %s", expected, r)
	}

	ValidPublishers = []string{
		"a",
		"b",
	}

	r = GetValidPublishersToString()
	expected = "The valid publishers are: a, b"

	if r != expected {
		t.Errorf("The provides are wrong. Expected: %s got %s", expected, r)
	}
}
