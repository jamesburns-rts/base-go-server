# Base build image
FROM golang:1.12-alpine AS build_base

ARG PROJECT_NAME
ARG BRANCH
ARG VERSION

LABEL projectName=$PROJECT_NAME
LABEL branch=$BRANCH
LABEL version=$VERSION
LABEL builder=true

# Install some dependencies needed to build the project
RUN apk add bash git gcc g++ libc-dev
WORKDIR /src

# Force the go compiler to use modules
ENV GO111MODULE=on

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

#This is the ‘magic’ step that will download all the dependencies that are specified in
# the go.mod and go.sum file.
# Because of how the layer caching system works in Docker, the  go mod download
# command will _ only_ be re-run when the go.mod or go.sum file change
# (or when we add another docker instruction this line)
RUN go mod download

FROM build_base AS server_builder
COPY . .
RUN /src/scripts/build_service.sh $VERSION $BRANCH

#In this last stage, we start from a fresh Alpine image, to reduce the image size and not ship the Go compiler in our production artifacts.
FROM alpine AS app
# # Finally we copy the statically compiled Go binary.
COPY --from=server_builder /src/bin/base-go-server /bin/base-go-server

EXPOSE 8080
CMD ["/bin/base-go-server"]
