// Package gopkg provides a minimal client to query the godoc.org JSON API.
//
// The newer pkg.go.dev site does not yet have an API, see
// https://github.com/golang/go/issues/36785.
// It's not clear how much longer the existing godoc.org API will keep working.
package gopkg

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client holds the settings to communicate with the godoc.org service.
// A zero Client is ready for use and provides default settings.
type Client struct {
	// The base URL of the service. If empty, uses gopkg.BaseURL.
	BaseURL string

	// The HTTP client to use for queries.  If nil, uses http.DefaultClient.
	HTTPClient
}

// get fetchs the specified URL and expects a successful reply with a JSON
// response body. In case of an error from the server, the decoded query result
// is returned along with a non-nil service error.
func (c Client) get(url string) (*queryResult, error) {
	var rsp *http.Response
	var err error
	if c.HTTPClient == nil {
		rsp, err = http.Get(url)
	} else {
		rsp, err = c.HTTPClient.Get(url)
	}
	if err != nil {
		return nil, fmt.Errorf("calling service: %w", err)
	}
	data, err := ioutil.ReadAll(rsp.Body)
	rsp.Body.Close()
	if err != nil {
		return nil, err
	}
	var result queryResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decoding result: %w", err)
	}
	if result.Error != nil {
		return &result, fmt.Errorf("service: %v", result.Error.Message)
	}
	return &result, nil
}

type term = [2]string

// makeURL assembles an API URL from the given method and terms.
// If provided, terms must consist of pairs of key=value strings.
func (c Client) makeURL(method string, terms ...term) string {
	var base string
	if c.BaseURL != "" {
		base = c.BaseURL + "/" + method
	} else {
		base = BaseURL + "/" + method
	}
	if len(terms) != 0 {
		query := make(url.Values)
		for _, t := range terms {
			query.Set(t[0], t[1])
		}
		return base + "?" + query.Encode()
	}
	return base
}

// Search lists the packages matching the specified query.  If the query
// succeeds but no packages match the query, Search returns nil, nil.
//
// API: /search?q=query
func (c Client) Search(_ context.Context, query string) ([]*Package, error) {
	res, err := c.get(c.makeURL("search", term{"q", query}))
	if err != nil {
		return nil, err
	}
	return res.Results, nil
}

// HTTPClient is the subset of the http.Client interface needed by this package.
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// BaseURL is the base URL for the JSON API endpoint, prior to addition of any
// path or query arguments.
const BaseURL = "https://api.godoc.org"

// Package describes the recorded information about a single package.
// Cloned from https://pkg.go.dev/github.com/golang/gddo/database#Package.
type Package struct {
	Name        string  `json:"name,omitempty"`
	Path        string  `json:"path"`
	ImportCount int     `json:"import_count"`
	Synopsis    string  `json:"synopsis,omitempty"`
	IsFork      bool    `json:"fork,omitempty"`
	NumStars    int     `json:"stars,omitempty"`
	Score       float64 `json:"score,omitempty"`

	// TODO: Import count does not appear to be indexed very consistently, most
	// of the packages report 0.
}

// queryResult is the top-level wrapper for responses from the JSON API.
type queryResult struct {
	Results []*Package  `json:"results"`
	Error   *queryError `json:"error"`
}

// queryError is the wrapper for error responses.
type queryError struct {
	Message string `json:"message"`
}
