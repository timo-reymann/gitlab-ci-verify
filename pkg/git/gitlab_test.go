package git

import (
	"reflect"
	"testing"
)

func TestParseGitlabRemoteUrlInfo(t *testing.T) {
	testCases := []struct {
		remoteUrl      string
		expectedInfo   *GitlabRemoteUrlInfo
		expectedErrMsg string
	}{
		{
			remoteUrl: "https://gitlab.com/timo-reymann/my-cool-project",
			expectedInfo: &GitlabRemoteUrlInfo{
				Hostname:       "gitlab.com",
				ClonedViaHttps: true,
				RepoSlug:       "timo-reymann/my-cool-project",
			},
		},
		{
			remoteUrl: "git@gitlab.com:timo-reymann/my-cool-project",
			expectedInfo: &GitlabRemoteUrlInfo{
				Hostname:       "gitlab.com",
				ClonedViaHttps: false,
				RepoSlug:       "timo-reymann/my-cool-project",
			},
		},
		{
			remoteUrl: "ssh://git@gitlab.com/timo-reymann/my-cool-project.git",
			expectedInfo: &GitlabRemoteUrlInfo{
				Hostname:       "gitlab.com",
				ClonedViaHttps: false,
				RepoSlug:       "timo-reymann/my-cool-project",
			},
		},
		{
			remoteUrl: "git://gitlab.com/timo-reymann/my-cool-project",
			expectedInfo: &GitlabRemoteUrlInfo{
				Hostname:       "gitlab.com",
				ClonedViaHttps: false,
				RepoSlug:       "timo-reymann/my-cool-project",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.remoteUrl, func(t *testing.T) {
			info, err := ParseGitlabRemoteUrlInfo(tc.remoteUrl)
			if err != nil {
				if tc.expectedErrMsg != "" {
					t.Fatalf("Expected no error, but got %v", err)
				} else if tc.expectedErrMsg != err.Error() {
					t.Fatalf("Expected error message %s, but got %s", tc.expectedErrMsg, err.Error())
				}
			}

			if !reflect.DeepEqual(tc.expectedInfo, info) {
				t.Fatalf("Expected %v, but got %v", tc.expectedInfo, info)
			}
		})
	}
}
