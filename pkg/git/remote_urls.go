package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/timo-reymann/gitlab-ci-verify/internal/logging"
	"slices"
)

func openRepository(path string) (*git.Repository, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func getRemoteUrls(r *git.Repository) ([]string, error) {
	var uniqueUrls = make(map[string]bool)

	remotes, err := r.Remotes()
	if err != nil {
		return nil, err
	}

	for _, remote := range remotes {
		for _, remoteUrl := range remote.Config().URLs {
			uniqueUrls[remoteUrl] = true
		}
	}

	var urls = make([]string, len(uniqueUrls))
	i := 0
	for url, _ := range uniqueUrls {
		logging.Debug("found git remote url", url)
		urls[i] = url
		i++
	}

	slices.Sort(urls)

	return urls, nil
}

// GetRemoteUrls for a git repository from the given path
func GetRemoteUrls(gitRepo string) ([]string, error) {
	repo, err := openRepository(gitRepo)
	if err != nil {
		return nil, err
	}

	return getRemoteUrls(repo)
}
