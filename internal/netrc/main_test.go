// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package netrc

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"
)

func TestReadUserNetrc(t *testing.T) {
	testCases := []struct {
		name          string
		netrcContent  string
		expectedLines []Line
		expectedErr   error
	}{
		{
			name:          "empty netrc file",
			netrcContent:  "",
			expectedLines: []Line{},
			expectedErr:   nil,
		},
		{
			name: "netrc file with valid entries",
			netrcContent: `
machine foo.com login bar password baz
`,
			expectedLines: []Line{
				{machine: "foo.com", login: "bar", password: "baz"},
			},
			expectedErr: nil,
		},
		{
			name: "netrc file with valid entries without login",
			netrcContent: `
machine foo.com password token
`,
			expectedLines: []Line{
				{machine: "foo.com", login: "", password: "token"},
			},
			expectedErr: nil,
		},
		{
			name: "netrc file with invalid entries",
			netrcContent: `
machine foo.com login bar
default login foo password bar
`,
			expectedLines: []Line{},
			expectedErr:   errors.New("missing password for machine foo.com"),
		},
		{
			name:          "error reading netrc file",
			netrcContent:  "",
			expectedLines: []Line{},
			expectedErr:   errors.New("mocked error reading netrc file"),
		},
		{
			name:          "no netrc file or env var",
			netrcContent:  "",
			expectedLines: []Line{},
			expectedErr:   errors.New("no such file or directory"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary netrc file
			tmpfile, err := ioutil.TempFile("", "netrc")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			// Write the netrc content to the file
			if _, err := tmpfile.Write([]byte(tc.netrcContent)); err != nil {
				t.Fatal(err)
			}

			// Set the NETRC environment variable to point to the temporary file
			origNetrc := os.Getenv("NETRC")
			os.Setenv("NETRC", tmpfile.Name())
			defer os.Setenv("NETRC", origNetrc)

			// Read the netrc file
			lines, err := ReadUserNetrc()

			// Check the results
			if err != nil {
				if !errors.Is(err, tc.expectedErr) {
					t.Errorf("expected error: %v, got: %v", tc.expectedErr, err)
				}
			} else if !compareLines(lines, tc.expectedLines) {
				t.Errorf("expected lines: %v, got: %v", tc.expectedLines, lines)
			}
		})
	}
}

func compareLines(lines1, lines2 []Line) bool {
	if len(lines1) != len(lines2) {
		return false
	}
	for i := range lines1 {
		if lines1[i].machine != lines2[i].machine ||
			lines1[i].login != lines2[i].login ||
			lines1[i].password != lines2[i].password {
			return false
		}
	}
	return true
}

func TestGetCredentials(t *testing.T) {
	testCases := []struct {
		name          string
		lines         []Line
		host          string
		expectedCreds *Credentials
		expectedErr   error
	}{
		{
			name:          "exact match",
			lines:         []Line{{machine: "example.com", login: "user", password: "password"}},
			host:          "example.com",
			expectedCreds: &Credentials{Login: "user", Password: "password"},
			expectedErr:   nil,
		},
		{
			name:          "match with scheme",
			lines:         []Line{{machine: "example.com", login: "user", password: "password"}},
			host:          "https://example.com",
			expectedCreds: &Credentials{Login: "user", Password: "password"},
			expectedErr:   nil,
		},
		{
			name:          "match with www prefix",
			lines:         []Line{{machine: "example.com", login: "user", password: "password"}},
			host:          "www.example.com",
			expectedCreds: &Credentials{Login: "user", Password: "password"},
			expectedErr:   nil,
		},
		{
			name:          "match without login prefix",
			lines:         []Line{{machine: "example.com", login: "", password: "token"}},
			host:          "example.com",
			expectedCreds: &Credentials{Login: "", Password: "token"},
			expectedErr:   nil,
		},
		{
			name:          "no match",
			lines:         []Line{{machine: "example.com", login: "user", password: "password"}},
			host:          "example.org",
			expectedCreds: nil,
			expectedErr:   ErrNoEntry,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			creds, err := GetCredentials(tc.lines, tc.host)
			if err != nil {
				if err != tc.expectedErr {
					t.Errorf("expected error: %v, got: %v", tc.expectedErr, err)
				}
				return
			}

			if !compareCredentials(creds, tc.expectedCreds) {
				t.Errorf("expected credentials: %v, got: %v", tc.expectedCreds, creds)
			}
		})
	}
}

func compareCredentials(c1, c2 *Credentials) bool {
	if c1 == nil && c2 == nil {
		return true
	}
	if c1 == nil || c2 == nil {
		return false
	}
	return c1.Login == c2.Login && c1.Password == c2.Password
}
