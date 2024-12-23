package buildinfo

import (
	"fmt"
	"runtime"
	"time"
)

var (
	// Version is the semantic version of the build
	Version = "dev"

	// GitCommit is the git SHA at build time
	GitCommit = "unknown"

	// BuildTime is the timestamp of the build
	BuildTime = "unknown"

	// BuildUser is the user who ran the build
	BuildUser = "unknown"
)

// Info represents build information
type Info struct {
	Version      string    `json:"version"`
	GitCommit    string    `json:"git_commit"`
	BuildTime    time.Time `json:"build_time"`
	BuildUser    string    `json:"build_user"`
	GoVersion    string    `json:"go_version"`
	Architecture string    `json:"architecture"`
	OS           string    `json:"os"`
}

// Get returns the current build information
func Get() *Info {
	buildTime, _ := time.Parse(time.RFC3339, BuildTime)
	return &Info{
		Version:      Version,
		GitCommit:    GitCommit,
		BuildTime:    buildTime,
		BuildUser:    BuildUser,
		GoVersion:    runtime.Version(),
		Architecture: runtime.GOARCH,
		OS:           runtime.GOOS,
	}
}

// String returns a formatted string of build information
func (i *Info) String() string {
	return fmt.Sprintf(
		"Version: %s\nGit Commit: %s\nBuild Time: %s\nBuild User: %s\nGo Version: %s\nArchitecture: %s\nOS: %s",
		i.Version,
		i.GitCommit,
		i.BuildTime.Format(time.RFC3339),
		i.BuildUser,
		i.GoVersion,
		i.Architecture,
		i.OS,
	)
}

// ShortVersion returns just the version string
func (i *Info) ShortVersion() string {
	return i.Version
}
