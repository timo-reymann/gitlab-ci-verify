package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// ErrInvalidAuthentication indicates the authentication is invalid
var ErrInvalidAuthentication = errors.New("invalid authentication")

// ErrNotFound indicates the resource or route does not exist
var ErrNotFound = errors.New("not found")

// Response as returned by gitlab api client methods
type Response struct {
	*http.Response
}

// CheckStatus for the response for common errors and return, according errors,
// if the HTTP status code does not indicate a general error nil is returned.
func (r *Response) CheckStatus() error {
	switch r.StatusCode {
	case http.StatusForbidden:
		return ErrInvalidAuthentication
	case http.StatusNotFound:
		return ErrNotFound
	}

	return nil
}

// UnmarshalJson from the response
func (r *Response) UnmarshalJson(val interface{}) error {
	raw, err := io.ReadAll(r.Response.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, &val)
}
