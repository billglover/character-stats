package skritter

import (
	"net/http"
	"net/url"
)

type Client struct {
	BaseURL    *url.URL
	UserAgent  string
	httpClient *http.Client
	token      string
}

func NewClient(token string) *Client {
	url, _ := url.Parse("https://legacy.skritter.com/api/v0")

	c := Client{
		BaseURL:    url,
		UserAgent:  "PoC App. Contact: bill@billglover.co.uk",
		httpClient: http.DefaultClient,
		token:      token,
	}

	return &c
}
