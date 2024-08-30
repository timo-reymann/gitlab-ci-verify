package ci_yaml

import (
	"context"
	"errors"
	formatconversion "github.com/timo-reymann/gitlab-ci-verify/pkg/format-conversion"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/git"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/gitlab/api"
	"time"
)

// ErrTimeout happens when no remote can validate the CI linting after the specified timeout
var ErrTimeout = errors.New("no remote url could validate")

// VerificationResultWithRemoteInfo contains both the remote that was used and the according lint result
type VerificationResultWithRemoteInfo struct {
	remoteInfo *git.GitlabRemoteUrlInfo
	lintResult *api.CiLintResult
}

// GetFirstValidationResult starts a parallel request to all remotes and tries to use the API to lint the CI file
// the first to report a result will be used as a validation result. If none of the remotes can produce a result or
// timeout is reached, the validation is aborted and the error is set accordingly.
func GetFirstValidationResult(remoteInfos []git.GitlabRemoteUrlInfo, baseUrlOverwrite string, file string, timeout time.Duration) (*VerificationResultWithRemoteInfo, error) {
	ciYaml, err := LoadRaw(file)
	if err != nil {
		return nil, err
	}

	ciJson, err := formatconversion.ToJson(ciYaml)
	if err != nil {
		return nil, err
	}

	result := make(chan VerificationResultWithRemoteInfo, 1)
	ctx, cancel := context.WithCancel(context.Background())
	for _, remoteInfo := range remoteInfos {
		go func(r *git.GitlabRemoteUrlInfo) {
			var baseUrl string
			if baseUrlOverwrite != "" {
				baseUrl = baseUrlOverwrite
			} else {
				baseUrl = r.Hostname
			}

			gl := api.NewClient(baseUrl, "tbd")
			lintRes, err := gl.LintCiYaml(ctx, r.RepoSlug, ciJson)
			if err != nil {
				return
			}

			// ignore sending on the closed channel
			defer func() {
				recover()
			}()

			result <- VerificationResultWithRemoteInfo{
				remoteInfo: r,
				lintResult: lintRes,
			}

			// close once the first request finishes
			close(result)
		}(&remoteInfo)
	}

	// wait for the first validation to arrive, or timeout
	var res VerificationResultWithRemoteInfo
	select {
	case res = <-result:
		break
	case <-time.After(timeout):
		err = ErrTimeout
	}

	// cancel context either way
	cancel()

	return &res, err
}
