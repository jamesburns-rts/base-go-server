package internal

import (
	"github.com/jamesburns-rts/base-go-server/internal/vm"
)

// variables from compile flags
var version, gitCommit, gitBranch, gitState, goVersion string

var Versions vm.Versions

func init() {
	Versions = vm.Versions{
		Version: version,
		Go:      goVersion,
		Git: vm.GitVersion{
			Commit: gitCommit,
			Branch: gitBranch,
			State:  gitState,
		},
	}
}
