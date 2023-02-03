package db

import (
	"github.com/jamesburns-rts/base-go-server/internal/call"
	"github.com/jamesburns-rts/base-go-server/internal/vm/page"
)

type UserExample struct {
	UserID      int    `db:"user_id"`
	ExampleName string `db:"example_name"`
}

func GetExamplesByUser(ctx call.Context, userID int, params page.Request) ([]UserExample, error) {
	examples := make([]UserExample, 0)
	err := ctx.DB().SelectContext(ctx, &examples, "select * from example.user_example where user_id = $1", userID)
	return examples, err
}
