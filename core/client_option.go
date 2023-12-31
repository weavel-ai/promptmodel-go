// This file was auto-generated by Fern from our API Definition.

package core

import (
	fmt "fmt"
	http "net/http"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		// Handle error if .env file is not found
	}
}

// ClientOption adapts the behavior of the generated client.
type ClientOption func(*ClientOptions)

// ClientOptions defines all of the possible client options.
// This type is primarily used by the generated code and is
// not meant to be used directly; use ClientOption instead.
type ClientOptions struct {
	BaseURL    string
	HTTPClient HTTPClient
	HTTPHeader http.Header
	ApiKey     string
}

// NewClientOptions returns a new *ClientOptions value.
// This function is primarily used by the generated code and is
// not meant to be used directly; use ClientOption instead.

func NewClientOptions() *ClientOptions {
	apiKey := os.Getenv("PROMPTMODEL_API_KEY")
	baseURL := os.Getenv("PROMPTMODEL_BACKEND_PUBLIC_URL")
	if baseURL == "" {
		baseURL = "https://promptmodel.up.railway.app"
	}

	return &ClientOptions{
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
		HTTPHeader: make(http.Header),
		ApiKey:     apiKey,
	}
}

// ToHeader maps the configured client options into a http.Header issued
// on every request.
func (c *ClientOptions) ToHeader() http.Header {
	header := c.cloneHeader()
	header.Set("Authorization", fmt.Sprintf("%v", c.ApiKey))
	return header
}

func (c *ClientOptions) cloneHeader() http.Header {
	return c.HTTPHeader.Clone()
}
