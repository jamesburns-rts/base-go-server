package main

import (
	"fmt"
	"github.com/jamesburns-rts/base-go-server/internal/config"
	"net/http"
	"os"
	"strconv"
)

// this mini-application is added to the docker image as the HEALTHCHECK
// it should be configured to run every 15 seconds or so

func main() {
	host := getLocalHostFromEnv(config.Defaults.LocalHost)
	port := getPortFromEnv(config.Defaults.Port)

	r, err := http.Get(fmt.Sprintf("http://%s:%d/api/public/healthz", host, port))
	if err != nil || r.StatusCode != 200 {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func getPortFromEnv(defaultPort int) int {
	portNumber, _ := strconv.Atoi(os.Getenv("APPLICATION_PORT"))
	if portNumber <= 0 {
		return defaultPort
	}
	return portNumber
}

func getLocalHostFromEnv(defaultHost string) string {
	localhost := os.Getenv("APPLICATION_LOCAL_HOST")
	if localhost == "" {
		return defaultHost
	}
	return localhost
}
