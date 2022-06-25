package jsonfile

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/ezeoleaf/larry/cache"
	"github.com/ezeoleaf/larry/config"
	"github.com/ezeoleaf/larry/domain"
	"github.com/go-redis/redis/v8"
)

func TestWriteFile(t *testing.T) {

	// content := []domain.Content{
	// 	{
	// 		Title:    StringToPointer("title1"),
	// 		Subtitle: StringToPointer("subtitle1"),
	// 		URL:      StringToPointer("url1"),
	// 		ExtraData: []string{
	// 			"extradata1",
	// 			"extradata2",
	// 		},
	// 	},
	// 	{
	// 		Title:    StringToPointer("title2"),
	// 		Subtitle: StringToPointer("subtitle2"),
	// 		URL:      StringToPointer("url2"),
	// 		ExtraData: []string{
	// 			"extradata3",
	// 			"extradata4",
	// 		},
	// 	},
	// }

	content := make([]domain.Content, 0)
	for i := 0; i < 10; i++ {
		content = append(content, domain.Content{
			Title:    StringToPointer(fmt.Sprintf("title-%d", i)),
			Subtitle: StringToPointer(fmt.Sprintf("subtitle-%d", i)),
			URL:      StringToPointer(fmt.Sprintf("url-%d", i)),
			ExtraData: []string{
				fmt.Sprintf("extradata-%d-1", i),
				fmt.Sprintf("extradata-%d-2", i),
			},
		})
	}

	f, err := os.Create("./temp.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.Encode(content)

	// f.WriteString(string(jsonString))
}

func TestReadFile(t *testing.T) {
	f, _ := os.Open("./temp.json")
	defer f.Close()

	decoder := json.NewDecoder(f)

	size := 1
	reservoir := make([]domain.Content, size)
	rand.Seed(time.Now().UnixNano())

	filteredData := make([]domain.Content, 0)

	// Read the array open bracket
	decoder.Token()

	count := 0
	for decoder.More() {
		data := new(domain.Content)
		decoder.Decode(data)

		// check cache/blacklist here for title

		filteredData = append(filteredData, *data)

		if count < size {
			reservoir[count] = *data
		} else {
			j := rand.Intn(count + 1)
			if j < size {
				fmt.Println("j:", j)
				reservoir[j] = *data
			}
		}
		count++
	}

	jsonString, _ := json.Marshal(reservoir)
	fmt.Println(string(jsonString))
}

func TestGetContentFromFile(t *testing.T) {
	mr, _ := miniredis.Run()
	ro := &redis.Options{
		Addr: mr.Addr(),
	}

	cfg := config.Config{
		LocalFile:  "./temp.json",
		FileFormat: "json",
	}

	cc := cache.NewClient(ro)
	p := Provider{Config: cfg, CacheClient: cc}
	p.getContentFromFile(cfg.LocalFile)

}

func TestGetContentFromReader(t *testing.T) {
	mr, _ := miniredis.Run()
	ro := &redis.Options{
		Addr: mr.Addr(),
	}

	cfg := config.Config{
		LocalFile:  "./temp.json",
		FileFormat: "json",
	}

	cc := cache.NewClient(ro)
	// TODO: load one in cache
	// TODO: load one in blacklist

	p := Provider{Config: cfg, CacheClient: cc}

	// strings.NewReader(tc.BlacklistFileContents)
	contentFile := `[{"Title":"title-0","Subtitle":"subtitle-0","URL":"url-0","ExtraData":["extradata-0-1","extradata-0-2"]},{"Title":"title-1","Subtitle":"subtitle-1","URL":"url-1","ExtraData":["extradata-1-1","extradata-1-2"]},{"Title":"title-2","Subtitle":"subtitle-2","URL":"url-2","ExtraData":["extradata-2-1","extradata-2-2"]}]`
	if content, err := p.getContentFromReader(strings.NewReader(contentFile)); err != nil {
		fmt.Println(err)
	} else {
		jsonString, _ := json.Marshal(content)
		fmt.Println(string(jsonString))
	}
}
