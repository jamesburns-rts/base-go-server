package util

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
)

// Version endpoint

type (
	versionResponse struct {
		Version string     `json:"version"`
		Go      string     `json:"go"`
		Echo    string     `json:"echo"`
		Git     gitVersion `json:"git"`
	}
	gitVersion struct {
		Commit string `json:"commit"`
		Branch string `json:"branch"`
		State  string `json:"state"`
	}
)

// variables from compile flags
var version, gitCommit, gitBranch, gitState, goVersion string
var Versions = versionResponse{
	Version: version,
	Go:      goVersion,
	Git: gitVersion{
		Commit: gitCommit,
		Branch: gitBranch,
		State:  gitState,
	},
}

func PrintVersionInfo() {

	bytes, err := json.MarshalIndent(Versions, "", " ")
	if err != nil {
		log.Warn("Unable to print version info")
		return
	}

	fmt.Println("Version Info:", string(bytes))
}
