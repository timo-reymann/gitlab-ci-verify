package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// ErrInvalidAuthentication indicates the authentication is invalid
var ErrInvalidAuthentication = errors.New("got HTTP/401 when trying to call lint API, verify your credentials are valid")

// ErrNotFound indicates the resource or route does not exist
var ErrNotFound = errors.New("got HTTP/404 when trying to call lint API, verify at least one remote is using GitLab")

// ErrForbidden indicates that the authentication is valid, but the requesting user does not have permission
var ErrForbidden = errors.New("got HTTP/403 when trying to call lint API, make sure you have access to the lint API for the project")

// Response as returned by gitlab api client methods
type Response struct {
	*http.Response
}

// CheckStatus for the response for common errors and return, according errors,
// if the HTTP status code does not indicate a general error nil is returned.
func (r *Response) CheckStatus() error {
	switch r.StatusCode {
	case http.StatusUnauthorized:
		return ErrInvalidAuthentication
	case http.StatusForbidden:
		return ErrForbidden
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
