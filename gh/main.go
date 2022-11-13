package gh

import (
	"fmt"
	"github.com/google/go-github/v48/github"
	"net/http"
	"os"
)

// Client is a wrapper around the github client
type Client struct {
	*github.Client
}

type Repo struct {
}

// NewClient creates a new github client
func NewClient() *Client {
	return &Client{
		Client: github.NewClient(nil),
	}
}

// makeHttpClient makes an HTTP client which can make authenticated requests to the GitHub API.
func makeHttpClient() (*http.Client, error) {
	token, err := getGithubTokenHeader()
	if err != nil {
		return nil, fmt.Errorf("getting token: %w", err)
	}
	t := &transport{
		token:               token,
		authHeader:          "Authorization",
		underlyingTransport: http.DefaultTransport,
	}

	c := &http.Client{
		Transport: t,
	}
	return c, nil
}

type transport struct {
	underlyingTransport http.RoundTripper
	authHeader          string
	token               string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add(t.authHeader, t.token)
	return t.underlyingTransport.RoundTrip(req)
}

func getGithubTokenHeader() (string, error) {
	val, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		return "", fmt.Errorf("GITHUB_TOKEN not set")
	}
	if val == "" {
		return "", fmt.Errorf("GITHUB_TOKEN is empty")
	}
	return fmt.Sprintf("Bearer %s", val), nil
}
