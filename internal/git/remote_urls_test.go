package git

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetRemoteUrls(t *testing.T) {
	testCases := []struct {
		repoPath       string
		expectedUrls   []string
		expectedErrMsg string
	}{
		{
			repoPath:     "multiple-remotes",
			expectedUrls: []string{"https://git.example.com", "https://git02.example.com"},
		},
		{
			repoPath:     "no-remotes",
			expectedUrls: []string{},
		},
		{
			repoPath:     "single-gitlab-com-remote",
			expectedUrls: []string{"https://gitlab.com/timo-reymann/foo.git"},
		},
		{
			repoPath:     "single-gitlab-selfhosted-remote",
			expectedUrls: []string{"https://gitlab.example.com/acme/frontend/ui-lib.git"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.repoPath, func(t *testing.T) {
			urls, err := GetRemoteUrls(fmt.Sprintf("testdata/repo-%s.git", tc.repoPath))
			if err != nil {
				if tc.expectedErrMsg != "" {
					t.Fatalf("Expected no error, but got %v", err)
				} else if tc.expectedErrMsg != err.Error() {
					t.Fatalf("Expected error message %s, but got %s", tc.expectedErrMsg, err.Error())
				}
			}

			if !reflect.DeepEqual(urls, tc.expectedUrls) {
				t.Fatalf("Expected urls to be %v, but got %v", tc.expectedUrls, urls)
			}
		})
	}
}
