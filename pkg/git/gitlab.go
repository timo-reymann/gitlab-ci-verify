package git

import (
	"github.com/chainguard-dev/git-urls"
	"strings"
)

// GitlabRemoteUrlInfo contains information about a gitlab remote
type GitlabRemoteUrlInfo struct {
	// Hostname of the Gitlab Instance
	Hostname string
	// ClonedViaHttps indicates the remote url is HTTPs
	ClonedViaHttps bool
	// RepoSlug Is the full qualified name of the project without trailing .git
	RepoSlug string
}

// ParseGitlabRemoteUrlInfo relevant for project detection for a given remote url string
func ParseGitlabRemoteUrlInfo(remoteUrl string) (*GitlabRemoteUrlInfo, error) {
	u, err := giturls.Parse(remoteUrl)
	if err != nil {
		return nil, err
	}

	return &GitlabRemoteUrlInfo{
		Hostname:       u.Hostname(),
		ClonedViaHttps: u.Scheme == "https",
		RepoSlug:       strings.TrimSuffix(strings.TrimPrefix(u.Path, "/"), ".git"),
	}, nil
}
