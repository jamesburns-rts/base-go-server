package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/jamesburns-rts/base-go-server/internal/api"
	"github.com/jamesburns-rts/base-go-server/internal/api/middleware"
	"github.com/jamesburns-rts/base-go-server/internal/clients/example"
	"github.com/jamesburns-rts/base-go-server/internal/config"
	"github.com/jamesburns-rts/base-go-server/internal/util"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
)

// http://patorjk.com/software/taag/#p=display&f=ANSI%20Shadow&t=GO-SERVER
func printBanner() {
	v := util.Versions
	fmt.Printf(`
...........................................................................
..██████╗..██████╗.......███████╗███████╗██████╗.██╗...██╗███████╗██████╗..
.██╔════╝.██╔═══██╗......██╔════╝██╔════╝██╔══██╗██║...██║██╔════╝██╔══██╗.
.██║..███╗██║...██║█████╗███████╗█████╗..██████╔╝██║...██║█████╗..██████╔╝.
.██║...██║██║...██║╚════╝╚════██║██╔══╝..██╔══██╗╚██╗.██╔╝██╔══╝..██╔══██╗.
.╚██████╔╝╚██████╔╝......███████║███████╗██║..██║.╚████╔╝.███████╗██║..██║.
..╚═════╝..╚═════╝.......╚══════╝╚══════╝╚═╝..╚═╝..╚═══╝..╚══════╝╚═╝..╚═╝.
...........................................................................
Application base-go-server %s

`, v.Version)

	versions := map[string]string{
		"Golang": v.Go,
		"Branch": v.Git.Branch,
		"Commit": v.Git.Commit,
		"Echo":   v.Echo,
	}
	for k, v := range versions {
		fmt.Printf("%s: %s\n", k, v)
	}
	println()
}

// Server small server struct to help with testing
// basically expands an echo server to have a collections of
// functions to call before shutting down
type Server struct {
	*echo.Echo
	closers []io.Closer
}

func main() {

	props, err := config.ReadProperties()
	if err != nil {
		log.Fatal("failed to read properties ", err)
	}
	config.PrintProperties(props)

	s, err := createServer(props)
	if err != nil {
		log.Fatal("unable to create server ", err)
	}

	printBanner()
	// start server
	go func() {
		log.Fatal(s.Start(fmt.Sprintf("%s:%d", props.LocalHost, props.Port)))
	}()

	s.waitAndShutdown()
}

// createServer Instantiates a full server without starting it
func createServer(props config.Application) (s *Server, err error) {

	// initial server setup
	s = &Server{Echo: echo.New()}

	// set up auth
	if err := middleware.ConfigureAuth(props); err != nil {
		return nil, errors.Wrap(err, "unable to configure auth")
	}

	// connect to database and migrate
	// - none here

	// tracing
	tracer := middleware.SetupTracer(props)
	s.closeOnShutdown(tracer)

	// clients
	exampleClient := example.NewClient(props)

	// configure server
	s.Echo.Use(
		echoMiddleware.Logger(),
		echoMiddleware.Recover(),
		middleware.ExampleClient(exampleClient),
		middleware.Tracer(tracer),
	)
	s.Echo.HTTPErrorHandler = api.HTTPErrorHandler
	s.Echo.Validator = util.Validator
	s.Echo.HideBanner = true
	s.Echo.Logger.SetLevel(props.EchoLogLevel())
	s.Echo.Debug = s.Logger.Level() == log.DEBUG
	log.SetLevel(props.EchoLogLevel())
	util.Versions.Echo = echo.Version

	api.AddRoutes(s.Echo)

	return s, nil
}

func (s *Server) waitAndShutdown() {

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	for _, f := range s.closers {
		if err := f.Close(); err != nil {
			log.Warn(err)
		}
	}

	log.Info("exiting")
}

func (s *Server) closeOnShutdown(closer io.Closer) {
	s.closers = append(s.closers, closer)
}
