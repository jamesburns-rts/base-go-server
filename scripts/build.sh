#!/bin/bash

cd `dirname $0`

VERSION=$(git tag | sort -V | tail -1)
GIT_COMMIT=$(git rev-list -1 HEAD)
GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
GIT_STATE=clean
if [[ $(git diff --stat) != '' ]]; then
    GIT_STATE=dirty
fi
GO_VERSION=$(go version | sed 's/ /-/g')

# Package where version info is kept
PKG="github.com/jamesburns-rts/base-go-server/internal" # /version.go

export CGO_ENABLED=0

mkdir -p ../bin
cd ../bin
go build -ldflags "-X $PKG.version=$VERSION -X $PKG.gitCommit=$GIT_COMMIT -X $PKG.gitBranch=$GIT_BRANCH -X $PKG.gitSummary=$GIT_SUMMARY -X $PKG.gitState=$GIT_STATE -X $PKG.goVersion=$GO_VERSION" ../cmd/base-go-server
