package git

import "github.com/go-git/go-git/v5"

func openRepository(path string) (*git.Repository, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func getRemoteUrls(r *git.Repository) ([]string, error) {
	urls := make([]string, 0)

	remotes, err := r.Remotes()
	if err != nil {
		return nil, err
	}

	for _, remote := range remotes {
		urls = append(urls, remote.Config().URLs...)
	}

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
