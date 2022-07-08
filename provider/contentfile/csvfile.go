package contentfile

import (
	"encoding/csv"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/ezeoleaf/larry/domain"
)

type CsvFileReader struct {
	skipHeader bool
}

func NewCsvFileReader(skipHeader bool) ContentFileReader {
	return CsvFileReader{skipHeader: skipHeader}
}

func (r CsvFileReader) getContentFromReader(handle io.Reader, skip func(string) bool) (*domain.Content, error) {
	size := 1
	var reservoir []string
	rand.Seed(time.Now().UnixNano())

	count := 0
	skipHeader := r.skipHeader
	csvReader := csv.NewReader(handle)
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// skip header line
		if skipHeader {
			skipHeader = false
			continue
		}

		if rec[0] == "" {
			log.Println("content missing title, skipping record")
			continue
		}

		if skip(rec[0]) {
			continue
		}

		// reservoir sampling technique
		if count < size {
			reservoir = rec
		} else {
			j := rand.Intn(count + 1)
			if j < size {
				reservoir = rec
			}
		}

		count++
	}

	if count > 0 {
		if content, err := convertCsvToContent(reservoir); err != nil {
			return nil, err
		} else {
			return content, nil
		}
	}

	return nil, nil
}

func convertCsvToContent(rec []string) (*domain.Content, error) {
	content := domain.Content{ExtraData: []string{}}
	if len(rec) > 0 {
		content.Title = StringToPointer(rec[0])
	}
	if len(rec) > 1 {
		content.Subtitle = StringToPointer(rec[1])
	}
	if len(rec) > 2 {
		content.URL = StringToPointer(rec[2])
	}
	if len(rec) > 3 {
		// number of extra data fields is variable for CSV
		content.ExtraData = make([]string, len(rec)-3)
		for i := 3; i < len(rec); i++ {
			content.ExtraData[i-3] = rec[i]
		}
	}
	return &content, nil
}
