package ci_yaml

import (
	"fmt"
	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/git"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

func mockCiValidate(status int, duration time.Duration, valid bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(duration)
		writer.WriteHeader(status)
		_, _ = writer.Write([]byte(fmt.Sprintf(`{ "valid": %s }`, strconv.FormatBool(valid))))
	}))
}

func TestGetFirstValidationResult(t *testing.T) {
	yamlContent, _ := os.ReadFile("test_data/valid.yml")

	res, err := GetFirstValidationResult(&ValidationResultInput{
		RemoteInfos: []git.GitlabRemoteUrlInfo{
			{
				Hostname:       mockCiValidate(http.StatusOK, 500*time.Millisecond, true).URL,
				ClonedViaHttps: true,
				RepoSlug:       "project/foo",
			},
			{
				Hostname:       mockCiValidate(http.StatusOK, 1*time.Second, true).URL,
				ClonedViaHttps: true,
				RepoSlug:       "project/foo",
			},
			{
				Hostname:       mockCiValidate(http.StatusOK, 2*time.Second, true).URL,
				ClonedViaHttps: true,
				RepoSlug:       "project/foo",
			},
		},
		Token:            "",
		BaseUrlOverwrite: "",
		CiYaml:           yamlContent,
		Timeout:          600 * time.Millisecond,
	})

	if err != nil {
		t.Fatal(err)
	}

	if res == nil {
		t.Fatal("No res given but did not throw error")
	}

}
