package rest

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type DoerFunc func(*http.Request) (*http.Response, error)

func (d DoerFunc) Do(req *http.Request) (*http.Response, error) {
	return d(req)
}

type NullCloser struct {
	io.Reader
}

type client struct {
	Doer
	baseURL string
}

func (c *client) BaseURL() string {
	return c.baseURL
}

func (n NullCloser) Close() error { return nil }

func TestClient(t *testing.T) {
	t.Run("BaseURL", func(t *testing.T) {
		t.Run("returns the base URL", func(t *testing.T) {
			c := &Client{baseURL: "http://example.com"}
			if got, want := c.BaseURL(), "http://example.com"; got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		})
	})
}

func Example() {
	// Create a client with a base URL.
	// this is a mock client, so more complex than your usual
	// Doer and BaseURLer setup.
	c := &client{
		DoerFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				// Return the URL as JSON string as the response body.
				Body: NullCloser{bytes.NewBuffer([]byte(fmt.Sprintf("%q", req.URL.String())))},
			}, nil
		}),
		"http://example.com"}

	// Defining an API method that makes a GET request and returns a string.
	path := APIMethod[Get, string]("/path/%s")

	// this is how you use it.
	// note that it returns a pointer to type T
	val, err := path.Do(c, "foo")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(*val)
	// Output: http://example.com/path/foo

}

func TestAPIMethod(t *testing.T) {

	t.Run("Do", func(t *testing.T) {
		t.Run("makes a request", func(t *testing.T) {
			receivedReq := &http.Request{}
			c := &client{DoerFunc(func(req *http.Request) (*http.Response, error) {
				receivedReq = req
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       NullCloser{bytes.NewBuffer([]byte(`"response"`))},
				}, nil
			}), "http://example.com"}
			m := APIMethod[Get, string]("/path")
			val, err := m.Do(c)
			if err != nil {
				t.Errorf("got %v, want nil", err)
			}
			if got, want := receivedReq.Method, "GET"; got != want {
				t.Errorf("got %q, want %q", got, want)
			}
			if got, want := receivedReq.URL.String(), "http://example.com/path"; got != want {
				t.Errorf("got %q, want %q", got, want)
			}
			if val == nil {
				t.Fatalf("got nil, want response")
			}
			if got, want := *val, "response"; got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		})
	})
}
