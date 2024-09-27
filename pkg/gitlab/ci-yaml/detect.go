package ci_yaml

import (
	"context"
	"errors"
	"github.com/timo-reymann/gitlab-ci-verify/internal/logging"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/git"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/gitlab/api"
	"time"
)

// ErrTimeout happens when no remote can validate the CI linting after the specified timeout
var ErrTimeout = errors.New("no remote url could validate due to a timeout")

// ErrNoResult happens when no remote can validate the CI as the requests don't succeed
var ErrNoResult = errors.New("no remote gitlab url could validate due to invalid responses, this indicates that either there is no remote defined which has an API endpoint or the API endpoint differs from the clone url")

// VerificationResultWithRemoteInfo contains both the remote that was used and the according lint result
type VerificationResultWithRemoteInfo struct {
	RemoteInfo *git.GitlabRemoteUrlInfo
	LintResult *api.CiLintResult
}

// GetFirstValidationResult starts a parallel request to all remotes and tries to use the API to lint the CI file
// the first to report a result will be used as a validation result. If none of the remotes can produce a result or
// timeout is reached, the validation is aborted and the error is set accordingly.
func GetFirstValidationResult(remoteInfos []git.GitlabRemoteUrlInfo, token string, baseUrlOverwrite string, ciYaml []byte, timeout time.Duration) (*VerificationResultWithRemoteInfo, error) {
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

			gl := api.NewClientWithMultiTokenSources(baseUrl, token)
			lintRes, err := gl.LintCiYaml(ctx, r.RepoSlug, ciYaml)
			if err != nil {
				logging.Warn("lint request failed for remote", remoteInfo.Hostname, ":", err.Error())
				return
			}

			// ignore sending on the closed channel
			defer func() {
				recover()
			}()

			result <- VerificationResultWithRemoteInfo{
				RemoteInfo: r,
				LintResult: lintRes,
			}
			logging.Verbose("lint request succeeded for remote", remoteInfo.Hostname, ": valid =", lintRes.Valid)

		}(&remoteInfo)
	}

	var err error

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

	if res.LintResult == nil {
		return nil, ErrNoResult
	}

	return &res, err
}
