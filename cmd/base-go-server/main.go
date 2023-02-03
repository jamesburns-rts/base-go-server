package main

import (
	"context"
	"fmt"
	"github.com/jamesburns-rts/base-go-server/internal"
	"github.com/jamesburns-rts/base-go-server/internal/db"
	"github.com/jamesburns-rts/base-go-server/internal/log"
	"io"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/jamesburns-rts/base-go-server/internal/api"
	"github.com/jamesburns-rts/base-go-server/internal/config"
	"github.com/labstack/echo/v4"
)

// http://patorjk.com/software/taag/#p=display&f=ANSI%20Shadow&t=GO-SERVER
func printBanner() {
	fmt.Printf(`
...........................................................................
..██████╗..██████╗.......███████╗███████╗██████╗.██╗...██╗███████╗██████╗..
.██╔════╝.██╔═══██╗......██╔════╝██╔════╝██╔══██╗██║...██║██╔════╝██╔══██╗.
.██║..███╗██║...██║█████╗███████╗█████╗..██████╔╝██║...██║█████╗..██████╔╝.
.██║...██║██║...██║╚════╝╚════██║██╔══╝..██╔══██╗╚██╗.██╔╝██╔══╝..██╔══██╗.
.╚██████╔╝╚██████╔╝......███████║███████╗██║..██║.╚████╔╝.███████╗██║..██║.
..╚═════╝..╚═════╝.......╚══════╝╚══════╝╚═╝..╚═╝..╚═══╝..╚══════╝╚═╝..╚═╝.
...........................................................................
`)

	versions := map[string]string{
		"Golang": internal.Versions.Go,
		"Branch": internal.Versions.Git.Branch,
		"Commit": internal.Versions.Git.Commit,
		"Echo":   internal.Versions.Echo,
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
	printBanner()

	props, err := config.ReadProperties()
	if err != nil {
		log.Fatal("failed to read properties ", err)
	}

	s, err := createServer(props)
	if err != nil {
		log.Fatal("unable to create server ", err)
	}

	// start server
	go func() {
		err := s.Start(fmt.Sprintf("%s:%d", props.LocalHost, props.Port))
		log.Fatal("start failed", err)
	}()

	s.waitAndShutdown()
}

// createServer Instantiates a full server without starting it
func createServer(props config.Application) (s *Server, err error) {

	// initial server setup
	s = &Server{Echo: echo.New()}

	// connect to database and migrate
	dbc, err := db.Connect(props.Database)
	if err != nil {
		return nil, err
	}
	s.closeOnShutdown(dbc)

	if err = db.Migrate(dbc.DB, props); err != nil {
		s.closeAllClosers()
		return nil, fmt.Errorf("error while migrating: %w", err)
	}

	middlewares, err := api.Middleware(props)
	if err != nil {
		s.closeAllClosers()
		return nil, err
	}

	// configure server
	s.Echo.Use(middlewares...)
	s.Echo.HTTPErrorHandler = api.HTTPErrorHandler
	s.Echo.HideBanner = true
	s.Echo.Debug = strings.ToUpper(props.LogLevel) == "DEBUG"
	internal.Versions.Echo = echo.Version

	ctxCreator := api.NewEchoContextCreator(props, dbc)
	s.addControllers(
		api.NewManagementController(props, ctxCreator),
		api.NewUsersController(props, ctxCreator),
	)

	return s, nil
}

type controller interface {
	AddRoutes(e *echo.Group)
}

func (s *Server) addControllers(controllers ...controller) {
	g := s.Echo.Group("")
	for _, c := range controllers {
		c.AddRoutes(g)
	}
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
		log.Error("failed to shutdown nicely", err)
	}

	s.closeAllClosers()

	log.Info("exiting")
}

func (s *Server) closeOnShutdown(closer io.Closer) {
	s.closers = append(s.closers, closer)
}

func (s *Server) closeAllClosers() {

	for _, f := range s.closers {
		if err := f.Close(); err != nil {
			log.Error("failed to close closer", err)
		}
	}
}
