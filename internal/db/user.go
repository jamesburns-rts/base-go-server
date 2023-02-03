package db

import (
	"database/sql"
	"github.com/jamesburns-rts/base-go-server/internal/call"
	"time"
)

type AppUser struct {
	ID           int        `db:"id"`
	Email        string     `db:"email"`
	PasswordHash string     `db:"password_hash"`
	FirstName    *string    `db:"first_name"`
	LastName     *string    `db:"last_name"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}

func FindUserByEmail(ctx call.Context, email string) (*AppUser, error) {
	var u AppUser
	err := ctx.DB().SelectContext(ctx, &u, "select * from example.app_user where email = $1", email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
