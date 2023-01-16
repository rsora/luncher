// Package user provides an example of a core business API. Right now these
// calls are just wrapping the data/store layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package user

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rsora/mywebapp/data/store/user"
	"go.uber.org/zap"
)

// Core manages the set of API's for user access.
type Core struct {
	log  *zap.SugaredLogger
	user user.Store
}

// NewCore constructs a core for user api access.
func NewCore(log *zap.SugaredLogger, db *sqlx.DB) Core {
	return Core{
		log:  log,
		user: user.NewStore(log, db),
	}
}

// Create inserts a new user into the database.
func (c Core) Create(ctx context.Context, nu user.NewUser, now time.Time) (user.User, error) {

	// PERFORM PRE BUSINESS OPERATIONS

	usr, err := c.user.Create(ctx, nu, now)
	if err != nil {
		return user.User{}, fmt.Errorf("create: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return usr, nil
}

// Update replaces a user document in the database.
func (c Core) Update(ctx context.Context, userID string, uu user.UpdateUser, now time.Time) error {

	// PERFORM PRE BUSINESS OPERATIONS

	if err := c.user.Update(ctx,
		userID, uu, now); err != nil {
		return fmt.Errorf("udpate: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return nil
}

// Delete removes a user from the database.
func (c Core) Delete(ctx context.Context, userID string) error {

	// PERFORM PRE BUSINESS OPERATIONS

	if err := c.user.Delete(ctx,
		userID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return nil
}

// Query retrieves a list of existing users from the database.
func (c Core) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]user.User, error) {

	// PERFORM PRE BUSINESS OPERATIONS

	users, err := c.user.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return users, nil
}

// QueryByID gets the specified user from the database.
func (c Core) QueryByID(ctx context.Context, userID string) (user.User, error) {

	// PERFORM PRE BUSINESS OPERATIONS

	usr, err := c.user.QueryByID(ctx,
		userID)
	if err != nil {
		return user.User{}, fmt.Errorf("query: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return usr, nil
}

// QueryByEmail gets the specified user from the database by email.
func (c Core) QueryByEmail(ctx context.Context, email string) (user.User, error) {

	// PERFORM PRE BUSINESS OPERATIONS

	usr, err := c.user.QueryByID(ctx,
		email)
	if err != nil {
		return user.User{}, fmt.Errorf("query: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return usr, nil
}

// Authenticate finds a user by their email and verifies their password. On
// success it returns a Claims User representing this user. The claims can be
// used to generate a token for future authentication.
func (c Core) Authenticate(ctx context.Context, now time.Time, email, password string) error {

	// PERFORM PRE BUSINESS OPERATIONS

	err := c.user.Authenticate(ctx, now, email, password)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	// PERFORM POST BUSINESS OPERATIONS

	return nil
}