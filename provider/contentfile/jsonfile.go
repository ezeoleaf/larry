package contentfile

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/ezeoleaf/larry/domain"
)

type JsonFileReader struct {
}

func NewJsonFileReader() ContentFileReader {
	return JsonFileReader{}
}

func (r JsonFileReader) getContentFromReader(handle io.Reader, skip func(string) bool) (*domain.Content, error) {
	size := 1
	reservoir := domain.Content{}
	rand.Seed(time.Now().UnixNano())

	decoder := json.NewDecoder(handle)
	if _, err := decoder.Token(); err != nil {
		if err.Error() == "EOF" {
			return nil, nil
		}
		return nil, err
	}

	count := 0
	for decoder.More() {
		data := new(domain.Content)
		if err := decoder.Decode(data); err != nil {
			return nil, err
		}

		if data.Title == nil || *data.Title == "" {
			log.Println("content missing title, skipping record")
			continue
		}

		if skip(*data.Title) {
			continue
		}

		// reservoir sampling technique
		if count < size {
			reservoir = *data
		} else {
			j := rand.Intn(count + 1)
			if j < size {
				reservoir = *data
			}
		}
		count++
	}

	if count > 0 {
		return &reservoir, nil
	}

	return nil, nil
}
