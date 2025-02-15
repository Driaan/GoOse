package goose

import (
	resty "github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type HtmlRequester interface {
	fetchHTML(string) (string, error)
}

// Crawler can fetch the target HTML page
type htmlrequester struct {
	config Configuration
}

// NewCrawler returns a crawler object initialised with the URL and the [optional] raw HTML body
func NewHtmlRequester(config Configuration) HtmlRequester {
	return htmlrequester{
		config: config,
	}
}

func (hr htmlrequester) fetchHTML(url string) (string, error) {
	client := resty.New()
	client.SetTimeout(hr.config.timeout)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36").
		Get(url)

	if err != nil {
		return "", errors.Wrap(err, "could not perform request on "+url)
	}
	if resp.IsError() {
		return "", &badRequest{Message: "could not perform request with " + url + " status code " + string(resp.StatusCode())}
	}
	return resp.String(), nil
}

type badRequest struct {
	Message string `json:"message,omitempty"`
}

func (BadRequest *badRequest) Error() string {
	return "Required request fields are not filled"
}
