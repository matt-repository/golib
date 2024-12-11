package version

import (
	"fmt"
	"log"
)

var (
	GitCommitId string
	GitTag      string
	BuildTime   string
)

func GetVersion() string {
	if GitTag == "" {
		return "unknown version"
	}
	return GitTag
}

func GetBuildTime() string {
	if BuildTime == "" {
		return "unknown build time"
	}
	return BuildTime
}

func GetGitCommitID() string {
	if GitCommitId == "" {
		return "unknown git commit_id"
	}
	return GitCommitId
}

func ShowVersion() {
	log.Printf("----------------Software Information---------------------\n")
	log.Printf(fmt.Sprintf("Version %s  \nBuild at %s \nGitInfo: %s\n", GitTag, BuildTime, GitCommitId))
	log.Printf("----------------------------------------------------------\n")
}
