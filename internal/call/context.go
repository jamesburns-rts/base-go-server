package call

import (
	"context"
	"database/sql"
	"github.com/jamesburns-rts/base-go-server/internal/log"
	"github.com/jmoiron/sqlx"
)

type (
	Context interface {
		context.Context
		ID() string
		UserID() int
		Log() *log.Logger
		DB() DBExecutor
		WrapInTransaction(f func() error) error
		WrapInReadOnlyTransaction(f func() error) error
	}

	DBExecutor interface {
		SelectContext(ctx context.Context, dest any, query string, args ...any) error
		NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)
		// maybe add some others in here
	}
)

type DefaultContext struct {
	context.Context
	ContextID    string // id of http call
	CallerUserID int
	Logger       *log.Logger

	DatabaseConn *sqlx.DB
	transaction  *sqlx.Tx
}

func (c *DefaultContext) ID() string {
	return c.ContextID
}

func (c *DefaultContext) UserID() int {
	return c.CallerUserID
}

func (c *DefaultContext) Log() *log.Logger {
	return c.Logger
}

func (c *DefaultContext) DB() DBExecutor {
	if c.transaction != nil {
		return c.transaction
	}
	return c.DatabaseConn
}

// WrapInTransaction wraps a function in a transaction, if that function returns
// an error, then the transaction is rolled back
func (c *DefaultContext) WrapInTransaction(f func() error) error {
	return c.wrapInTransaction(false, f)
}

// WrapInReadOnlyTransaction wraps a function in a transaction, if that function returns
// an error, then the transaction is rolled back
func (c *DefaultContext) WrapInReadOnlyTransaction(f func() error) error {
	return c.wrapInTransaction(true, f)
}

func (c *DefaultContext) wrapInTransaction(readOnly bool, f func() error) error {
	if c.transaction == nil {
		var err error
		c.transaction, err = c.DatabaseConn.BeginTxx(c, &sql.TxOptions{ReadOnly: readOnly})
		if err != nil {
			return err
		}
	}

	silentlyRollbackTransaction := func() {
		err := c.transaction.Rollback()
		c.transaction = nil
		if err != nil {
			c.Log().Error("failed to rollback transaction", err)
		}
	}

	defer func() {
		// catch panics
		if r := recover(); r != nil {
			silentlyRollbackTransaction()
			panic(r)
		}
	}()

	// execute wrapped function
	if err := f(); err != nil {
		silentlyRollbackTransaction()
		return err
	}

	err := c.transaction.Commit()
	c.transaction = nil
	return err
}

func LoggerFromContext(ctx context.Context) *log.Logger {
	if callContext, ok := ctx.(Context); ok {
		return callContext.Log()
	}
	return log.DefaultLogger
}
