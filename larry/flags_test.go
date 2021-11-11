package larry

import (
	"testing"

	cfg "github.com/ezeoleaf/larry/config"
)

func TestGetFlags(t *testing.T) {
	mockConfig := cfg.Config{}
	flags := GetFlags(&mockConfig)
	if flags == nil {
		t.Errorf("Expected flags, got %s", flags)
	}

	if len(flags) < 1 {
		t.Error("No flags found")
	}
}
