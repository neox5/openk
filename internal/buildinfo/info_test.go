package buildinfo_test

import (
	"testing"
	"time"

	"github.com/neox5/openk/internal/buildinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInfo_Get(t *testing.T) {
	// Save original values
	origVersion := buildinfo.Version
	origGitCommit := buildinfo.GitCommit
	origBuildTime := buildinfo.BuildTime
	origBuildUser := buildinfo.BuildUser

	// Restore after test
	defer func() {
		buildinfo.Version = origVersion
		buildinfo.GitCommit = origGitCommit
		buildinfo.BuildTime = origBuildTime
		buildinfo.BuildUser = origBuildUser
	}()

	t.Run("success cases", func(t *testing.T) {
		t.Run("returns info with test values", func(t *testing.T) {
			// Setup test values
			buildinfo.Version = "v1.0.0"
			buildinfo.GitCommit = "abc123"
			buildinfo.BuildTime = "2024-01-01T12:00:00Z"
			buildinfo.BuildUser = "test-user"

			info := buildinfo.Get()
			require.NotNil(t, info)

			assert.Equal(t, "v1.0.0", info.Version)
			assert.Equal(t, "abc123", info.GitCommit)
			assert.Equal(t, "test-user", info.BuildUser)

			expectedTime, err := time.Parse(time.RFC3339, "2024-01-01T12:00:00Z")
			require.NoError(t, err)
			assert.Equal(t, expectedTime, info.BuildTime)
		})
	})

	t.Run("handles invalid build time", func(t *testing.T) {
		buildinfo.BuildTime = "invalid"
		info := buildinfo.Get()
		assert.NotNil(t, info)
		assert.True(t, info.BuildTime.IsZero())
	})
}

func TestInfo_String(t *testing.T) {
	t.Run("formats info correctly", func(t *testing.T) {
		buildTime, err := time.Parse(time.RFC3339, "2024-01-01T12:00:00Z")
		require.NoError(t, err)

		info := &buildinfo.Info{
			Version:      "v1.0.0",
			GitCommit:    "abc123",
			BuildTime:    buildTime,
			BuildUser:    "test-user",
			GoVersion:    "go1.21.0",
			Architecture: "amd64",
			OS:           "linux",
		}

		expected := "Version: v1.0.0\n" +
			"Git Commit: abc123\n" +
			"Build Time: 2024-01-01T12:00:00Z\n" +
			"Build User: test-user\n" +
			"Go Version: go1.21.0\n" +
			"Architecture: amd64\n" +
			"OS: linux"

		assert.Equal(t, expected, info.String())
	})
}

func TestInfo_ShortVersion(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		tests := []struct {
			name     string
			version  string
			expected string
		}{
			{
				name:     "returns release version",
				version:  "v1.2.3",
				expected: "v1.2.3",
			},
			{
				name:     "returns dev version",
				version:  "dev",
				expected: "dev",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				info := &buildinfo.Info{Version: tt.version}
				assert.Equal(t, tt.expected, info.ShortVersion())
			})
		}
	})
}
