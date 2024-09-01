// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netrc

import (
	"errors"
	"github.com/timo-reymann/gitlab-ci-verify/internal/logging"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ErrNoEntry represents the case that there is no entry for a given hostname in the netrc
var ErrNoEntry = errors.New("no entry for hostname")

// Line representing a netrc entry
type Line struct {
	machine  string
	login    string
	password string
}

func parseNetrc(data string) []Line {
	// See https://www.gnu.org/software/inetutils/manual/html_node/The-_002enetrc-file.html
	// for documentation on the .netrc format.
	var nrc []Line
	var l Line
	inMacro := false
	for _, line := range strings.Split(data, "\n") {
		if inMacro {
			if line == "" {
				inMacro = false
			}
			continue
		}

		f := strings.Fields(line)
		i := 0
		for ; i < len(f)-1; i += 2 {
			// Reset at each "machine" token.
			// “The auto-login process searches the .netrc file for a machine token
			// that matches […]. Once a match is made, the subsequent .netrc tokens
			// are processed, stopping when the end of file is reached or another
			// machine or a default token is encountered.”
			switch f[i] {
			case "machine":
				l = Line{machine: f[i+1]}
			case "default":
				break
			case "login":
				l.login = f[i+1]
			case "password":
				l.password = f[i+1]
			case "macdef":
				// “A macro is defined with the specified name; its contents begin with
				// the next .netrc line and continue until a null line (consecutive
				// new-line characters) is encountered.”
				inMacro = true
			}
			if l.machine != "" && l.login != "" && l.password != "" {
				nrc = append(nrc, l)
				l = Line{}
			}
		}

		if i < len(f) && f[i] == "default" {
			// “There can be only one default token, and it must be after all machine tokens.”
			break
		}
	}

	return nrc
}

func netrcPath() (string, error) {
	if env := os.Getenv("NETRC"); env != "" {
		return env, nil
	}
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	base := ".netrc"
	if runtime.GOOS == "windows" {
		base = "_netrc"
	}
	return filepath.Join(dir, base), nil
}

// ReadUserNetrc looks up the netrc for the user if set and returns the parsed result
// in case there is no netrc an error is returned
func ReadUserNetrc() ([]Line, error) {
	path, err := netrcPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	netrc := parseNetrc(string(data))
	logging.Debug("Found netrc for user in", path)
	return netrc, nil
}

// Credentials for a host
type Credentials struct {
	// Login is typically the username
	Login string
	// Password contains the secret to authenticate with
	Password string
}

// GetCredentials from given netrc lines and the given host
func GetCredentials(lines []Line, host string) (*Credentials, error) {
	hostWithoutScheme := strings.TrimPrefix(strings.TrimPrefix(host, "https://"), "http://")

	for _, line := range lines {
		if line.machine == hostWithoutScheme || strings.HasSuffix(hostWithoutScheme, line.machine) {
			logging.Verbose("Found credentials for", host, "in netrc")
			return &Credentials{
				Login:    line.login,
				Password: line.password,
			}, nil
		}
	}

	return nil, ErrNoEntry
}
