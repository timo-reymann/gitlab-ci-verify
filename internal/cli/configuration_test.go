package cli

import (
	"os"
	"testing"
)

func TestGitlabBaseUrlOverwrite(t *testing.T) {
	tests := []struct {
		name          string
		gitlabBaseUrl string
		expected      string
	}{
		{
			name:          "auto-detect value",
			gitlabBaseUrl: AutoDetectValue,
			expected:      "",
		},
		{
			name:          "explicit value",
			gitlabBaseUrl: "https://gitlab.example.com",
			expected:      "https://gitlab.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &Configuration{GitlabBaseUrl: tt.gitlabBaseUrl}
			if got := conf.GitlabBaseUrlOverwrite(); got != tt.expected {
				t.Errorf("GitlabBaseUrlOverwrite() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsCIEnv(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected bool
	}{
		{
			name:     "CI environment variable set",
			envValue: "true",
			expected: true,
		},
		{
			name:     "CI environment variable not set",
			envValue: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("CI", tt.envValue)
			} else {
				os.Unsetenv("CI")
			}

			conf := &Configuration{}
			if got := conf.IsCIEnv(); got != tt.expected {
				t.Errorf("IsCIEnv() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedErr error
	}{
		{
			name:        "help flag",
			args:        []string{"--help"},
			expectedErr: ErrAbort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = append([]string{"cmd"}, tt.args...)
			conf := NewConfiguration()
			err := conf.Parse("")
			if err != tt.expectedErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.expectedErr)
			}
		})
	}
}
