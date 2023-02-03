package db

import (
	"fmt"
	"github.com/jamesburns-rts/base-go-server/internal/config"
	"github.com/jmoiron/sqlx"
	"strings"
)

func Connect(props config.Database) (*sqlx.DB, error) {
	params := []string{
		fmt.Sprintf("host=%s", props.Host),
		fmt.Sprintf("port=%d", props.Port),
		fmt.Sprintf("dbname=%s", props.Database),
		fmt.Sprintf("user=%s", props.Username),
		fmt.Sprintf("password=%s", props.Password),
	}
	if props.SSLDisable {
		fmt.Println("appending ssl disable")
		params = append(params, "sslmode=disable")
	}

	fmt.Println("conniection", strings.Join(params, " "))
	return sqlx.Connect("postgres", strings.Join(params, " "))
}
