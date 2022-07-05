package oauth2

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	oauth "golang.org/x/oauth2"
)

const (
	TWITTER_AUTH_URL  = "https://twitter.com/i/oauth2/authorize"
	TWITTER_TOKEN_URL = "https://api.twitter.com/2/oauth2/token"

	ENDPOINT = "https://api.twitter.com/2/"

	TWEET_ENDPOINT = ENDPOINT + "tweets"
)

type Twitter struct {
	client *http.Client
}

var filename string

func init() {
	filename = getFileName("larry.env")
}

// Returns a new oauth2 config with the given client id and secret
func NewConfig(id, secret string) *oauth.Config {
	return &oauth.Config{
		ClientID:     id,
		ClientSecret: secret,
		Scopes:       []string{"tweet.write", "tweet.read", "users.read", "offline.access"},
		Endpoint: oauth.Endpoint{
			AuthURL:   TWITTER_AUTH_URL,
			TokenURL:  TWITTER_TOKEN_URL,
			AuthStyle: oauth.AuthStyleAutoDetect,
		},
		RedirectURL: "http://localhost:8080/callback",
	}
}

// Returns a new http Client authorized via oauth2 flow
func NewClient(ctx context.Context, conf *oauth.Config) (*Twitter, error) {
	var err error

	// getting initial values of the oauth2 token
	tok := getToken()

	// no need for authorization, if AccessToken and RefreshToken are found
	if tok.AccessToken == "" && tok.RefreshToken == "" {
		state := "solid"

		// set additional required parameter for authorization
		code_challenge := oauth.SetAuthURLParam("code_challenge", "larry")
		code_challenge_method := oauth.SetAuthURLParam("code_challenge_method", "plain")

		// generating the auth url
		url := conf.AuthCodeURL(state, oauth.AccessTypeOffline, code_challenge, code_challenge_method)

		err = openUrl(url)
		if err != nil {
			return nil, err
		}

		cd := make(chan string)
		st := make(chan string)
		// handlecallback in the background
		go handleCallback(cd, st)

		var callbackState, code string
		// wait for the callback, to get state and code
		select {
		case callbackState = <-st:
			code = <-cd
		case <-time.After(time.Minute):
			return nil, errors.New("err: authorization failed[timed out!!]")
		}

		// verify the state received from the callback url
		if callbackState == state {
			codeVerifier := oauth.SetAuthURLParam("code_verifier", "larry")
			clientId := oauth.SetAuthURLParam("client_id", conf.ClientID)

			tok, err = conf.Exchange(ctx, code, codeVerifier, clientId)
			if err != nil {
				return nil, errors.New("err: couldn't generate token")
			}

			// store the token for future use, to avoid re-authorization
			err = storeToken(tok)
			if err != nil {
				return nil, errors.New("err: failed to store tokens")
			}
		} else {
			return nil, errors.New("err: state verification failed")
		}
	} else {
		// regenerate token if expired
		expired := tok.Expiry.Round(0).Add(-(10 * time.Second)).Before(time.Now())
		if expired {
			log.Println("Re-generating token...")
			tok, err = regenerateToken(ctx, conf, tok)
			if err != nil {
				return nil, errors.New("err: failed to regenerate tokens")
			}

			// Update the tokens
			err = storeToken(tok)
			if err != nil {
				return nil, errors.New("err: failed to update tokens")
			}
		}
	}

	httpClient := conf.Client(ctx, tok)

	return &Twitter{client: httpClient}, nil
}

// Handles callback url during authorization, takes two channels for code and state.
//
// Fetches code and state from the redirect url, and push it to the respective channels
func handleCallback(cd chan<- string, st chan<- string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		state := r.FormValue("state")
		if state != "" && code != "" {
			st <- state
			cd <- code
		} else {
			log.Fatal("err: failed to parse parameters")
		}
	})

	log.Println("Waiting for app authorization...")
	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}

// Sends tweet on behalf of the authenticated user,
// accepts authenticated httpClient and tweet message.
// Returns ID of newly created tweet on success.
func (t *Twitter) Update(tweet string) (id string, err error) {
	req, err := json.Marshal(map[string]string{"text": tweet})
	if err != nil {
		return "", err
	}

	resp, err := t.client.Post(TWEET_ENDPOINT, "application/json",
		bytes.NewBuffer(req))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res map[string]map[string]string
	json.NewDecoder(resp.Body).Decode(&res)

	body := "id: " + res["data"]["id"] + "\n"
	return body, nil
}

// Opens url in the system default browser
func openUrl(url string) error {
	var err error
	if url != "" {
		switch runtime.GOOS {
		case "linux":
			err = exec.Command("xdg-open", url).Start()
		case "windows":
			err = exec.Command("rundll32", "url.dll.FileProtocolHandler", url).Start()
		case "darwin":
			err = exec.Command("open", url).Start()
		default:
			log.Fatal("Unsupported platform")
		}

		return err
	} else {
		return errors.New("err: url is empty")
	}
}
