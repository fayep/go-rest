// A package for making requests to a Rest API.
// This package avoids making you use particular
// Underlying clients.  As long as it implements
// the Doer interface, you can use it with this package.
// The package also allows you to define the base URL
// for the API, and then define methods that can be
// formatted with input to create a URL to be added
// to the base URL.

// The package is eminently testable, as it uses interfaces
// for the Doer and BaseURLer, and the client struct is
// easily mockable.

// The package is also easily extensible, as you can define
// new methods by defining new types that implement the
// Requester interface, and then defining new APIMethods
// with the new Requester type.

package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Doer is an interface that can make HTTP requests.
// It is implemented by http.Client for example.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// BaseURLer is an interface that returns the base URL for the API.
type BaseURLer interface {
	BaseURL() string
}

type RoundTripFunc func(*http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type Client struct {
	*http.Client
	baseURL string
}

// Ensure Client implements BaseURLer
func (c *Client) BaseURL() string {
	return c.baseURL
}

// NewBearerTokenClient creates a new client with a bearer token.
// The token is added to the Authorization header on each request.
// You can just use an http.Client if you don't need this.
// The baseURL is the base URL for the API.
func NewBearerTokenClient(baseURL, token string) *Client {
	httpClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			req.Header.Set("Authorization", "Bearer "+token)
			return http.DefaultTransport.RoundTrip(req)
		}),
	}
	return &Client{Client: httpClient, baseURL: baseURL}
}

func makeRequest(c Doer, method string, url string, input ...any) (*http.Response, error) {
	var body io.Reader
	var header http.Header
	for _, i := range input {
		switch v := i.(type) {
		case io.Reader:
			body = v
		case http.Header:
			header = v
		}
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header = header
	return c.Do(req)
}

// Requester is an interface that defines the method
// to be used in the request.
// It needs to be understood by your Doer.
type Requester interface {
	Method() string
}

// Get, Post, and Put are types that implement the Requester interface.
// They are used to define the method to be used in the request.
// They are entirely static to simplify the APIMethod type.
type Get struct{}
type Post struct{}
type Put struct{}

// static methods for interface
func (m Get) Method() string {
	return "GET"
}

func (m Post) Method() string {
	return "POST"
}

func (m Put) Method() string {
	return "PUT"
}

// APIMethod is a format string describing the path
// to be added to the base URL, any parameters to Do() will be
// formatted into the string.
type APIMethod[R Requester, T any] string

// Do makes a request to the API with the given input.
// The input is formatted into the path string.
// The response is decoded into the output.
// The output is a pointer to the type T.
// If the Doer is also a BaseURLer,
// the base URL is prepended to the path.
func (a APIMethod[R, T]) Do(c Doer, input ...any) (*T, error) {
	var method R
	baseURL := ""
	if b, ok := c.(BaseURLer); ok {
		baseURL = b.BaseURL()
	}
	output := new(T)
	args, others := splitArgs(string(a), input...)
	resp, err := makeRequest(
		c,
		method.Method(),
		fmt.Sprintf(baseURL+string(a), args...),
		others..., // is there header/body handling required?
	)
	if err != nil {
		return nil, err
	}
	err = decodeResponse(resp, output)
	return output, err
}

func decodeResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}
