package oauth2

import (
	"bufio"
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	oauth "golang.org/x/oauth2"
)

func getFileName(name string) string {
    home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// file will be stored in home directory of current user
	home += "/." + name

    return home
}

// Regenerates token for the given client.
//
// Returns new token on success, else old token with error.
func regenerateToken(ctx context.Context, conf *oauth.Config, tok *oauth.Token) (*oauth.Token, error) {
	var err error

	if tok != nil && conf != nil {
		source := conf.TokenSource(ctx, tok)
		tok, err = source.Token()
		if err != nil {
			return tok, err
		}
	} else {
		return nil, errors.New("err: token or config is invalid")
	}

	return tok, err
}

// Stores token inside env file for future use
func storeToken(tok *oauth.Token) error {
	if tok != nil {
		store("TWITTER_ACCESS_TOKEN", tok.AccessToken)
		store("TWITTER_REFRESH_TOKEN", tok.RefreshToken)
		store("TWITTER_TOKEN_EXPIRY", tok.Expiry.Format(time.RFC1123))

		return nil
	} else {
		return errors.New("err: token is invalid")
	}
}

// Retrieves token information of previous session from the env file,
// returns empty token if not found.
func getToken() *oauth.Token {
	var (
		err error
		accessToken, refreshToken string
		expire                    time.Time
	)

	// get the values of token from env file
	if t := get("TWITTER_TOKEN_EXPIRY"); t != "" {
		expire, err = time.Parse(time.RFC1123, t)
		if err == nil {
			accessToken = get("TWITTER_ACCESS_TOKEN")
			refreshToken = get("TWITTER_REFRESH_TOKEN")
		}
	}

	return &oauth.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expiry:       expire,
	}
}

// Stores key-value pair in target file.
//
// Returns error, if any
func store(key string, value string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	content, err := getFileContent(f)
	if err != nil {
		return err
	}

	var data string
	// if file is not empty then copy whole content and then make the changes
	if len(content) > 0 {
		found, index := getKey(content, key)

		// if not found then just append to file
		if !found {
			data = strings.Join(content, "\n") + "\n"
			data += key + "=" + value + "\n"
		} else {
			temp := strings.Split(content[index], "=")
			temp[1] = value
			content[index] = strings.Join(temp, "=")
			data = strings.Join(content, "\n") + "\n"
		}

	} else {
		data = key + "=" + value + "\n"
	}

	// overwrite the file content
	err = f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(data))
	if err != nil {
		return err
	}
	f.Sync()

	return nil
}

// Returns value of the provided key(if present)
func get(key string) (value string) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return ""
	}
	defer f.Close()

	content, err := getFileContent(f)
	if err != nil {
		return ""
	}

	if ok, i := getKey(content, key); ok {
		value = strings.Split(content[i], "=")[1]
	}

	return value
}

// Searches for the given key inside an env file,
// returns found status with the position.
//
// If not found, position returned is -1
func getKey(content []string, key string) (found bool, index int) {
	for k, v := range content {
		// ignore comments and empty lines
		if !strings.HasPrefix(strings.TrimSpace(v), "#") && len(v) != 0 {
			if strings.Index(v, "=") != -1 && strings.Contains(v, key) {
				return true, k
			}
		}
	}

	return false, -1
}

// Get file's whole content in form of string array,
// each element representing a line from the file.
func getFileContent(file *os.File) ([]string, error) {
	var lines []string

	if file != nil {
		scan := bufio.NewScanner(file)
		for scan.Scan() {
			lines = append(lines, scan.Text())
		}

		if err := scan.Err(); err != nil {
			return lines, err
		}
	}

	return lines, nil
}
