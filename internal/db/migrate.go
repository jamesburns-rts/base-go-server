package db

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/jamesburns-rts/base-go-server/internal/config"
	"github.com/jamesburns-rts/base-go-server/internal/log"
	"github.com/jamesburns-rts/base-go-server/internal/util"
	"github.com/jamesburns-rts/base-go-server/internal/util/ptr"
	"github.com/jmoiron/sqlx"
	"os"
	"strings"
	// required for goose
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embeddedMigrations embed.FS

// Migrate the database using the scripts in "migrations"
func Migrate(dbc *sql.DB, props config.Application) error {
	_ = goose.SetDialect("postgres")
	goose.SetBaseFS(embeddedMigrations)
	goose.SetLogger(gooseLogger{logger: log.NewLogger(props.LogLevel, "goose")})

	goose.AddNamedMigrationNoTx("20230203095700_initial_user.go", func(d *sql.DB) error {
		return initialUserMigration(props, d)
	}, nil)

	if err := goose.Up(dbc, "migrations", goose.WithAllowMissing()); err != nil {
		return fmt.Errorf("during goose up: %w", err)
	}
	return nil
}

func initialUserMigration(props config.Application, d *sql.DB) error {
	password := props.InitialPassword
	if password == "" {
		var err error
		password, err = util.RandomAlphanumericString(40)
		if err != nil {
			return fmt.Errorf("failed to make password: %w", err)
		}
	}
	sd := sqlx.NewDb(d, "postgres")
	_, err := sd.NamedExec(
		`insert into example.app_user (email, password_hash, first_name, last_name) 
values (:email, :password_hash, :first_name, :last_name)`,
		AppUser{
			Email:        "admin@test.com",
			PasswordHash: "",
			FirstName:    ptr.To("Initial"),
			LastName:     ptr.To("User"),
		})

	log.Info("Created initial user with password", "password", password)
	return err
}

type gooseLogger struct {
	logger *log.Logger
}

func (g gooseLogger) Fatal(v ...interface{}) {
	msg := fmt.Sprint(v...)
	g.logger.Error(msg, errors.New(msg))
	os.Exit(1)
}

func (g gooseLogger) Fatalf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	g.logger.Error(msg, errors.New(msg))
	os.Exit(1)
}

func (g gooseLogger) Print(v ...interface{}) {
	msg := fmt.Sprint(v...)
	g.logger.Info(msg)
}

func (g gooseLogger) Println(v ...interface{}) {
	msg := fmt.Sprint(v...)
	g.logger.Info(msg)
}

func (g gooseLogger) Printf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	msg = strings.Trim(msg, "\n")
	g.logger.Info(msg)
}
