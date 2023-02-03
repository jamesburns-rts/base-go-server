# Base build image
FROM golang:1.19-alpine AS build_base

# Install some dependencies needed to build the project
RUN apk add bash git ca-certificates openssh
WORKDIR /src

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
RUN /src/scripts/build.sh
RUN CGO_ENABLED=0 go build -installsuffix 'static' -o /src/bin/healthcheck /src/cmd/healthcheck

#In this last stage, we start from a fresh Alpine image, to reduce the image size and not ship the Go compiler in our production artifacts.
FROM scratch AS app
# # Finally we copy the statically compiled Go binary.
COPY --from=server_builder /src/bin/base-go-server /bin/base-go-server
COPY --from=server_builder /src/bin/healthcheck /bin/healthcheck

# This allows us to make https calls from server
COPY --from=server_builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

HEALTHCHECK --interval=15s --timeout=5s --start-period=5s --retries=3 CMD ["/bin/healthcheck"]

EXPOSE 8080
ENV APP_LOCAL_HOST=0.0.0.0
CMD ["/bin/base-go-server"]
