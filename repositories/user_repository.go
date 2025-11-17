package repositories

import (
	"context"
	"database/sql"
	"time"

	dbSqlc "github.com/DraconDev/go-templ-htmx-ex/db/sqlc"
	"github.com/google/uuid"
)

// User represents a user in the application
type User struct {
	ID        uuid.UUID `json:"id"`
	AuthID    string    `json:"auth_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Picture   string    `json:"picture"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository handles user data access operations
type UserRepository struct {
	queries *dbSqlc.Queries
}

// NewUserRepository creates a new user repository
func NewUserRepository(queries *dbSqlc.Queries) *UserRepository {
	return &UserRepository{
		queries: queries,
	}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(ctx context.Context, user *User) (*User, error) {
	if r.queries == nil {
		return nil, nil // Database not connected
	}

	picture := sql.NullString{String: user.Picture, Valid: user.Picture != ""}
	isAdmin := sql.NullBool{Bool: user.IsAdmin, Valid: true}

	dbUser, err := r.queries.CreateUser(ctx, dbSqlc.CreateUserParams{
		AuthID:   user.AuthID,
		Email:    user.Email,
		Name:     user.Name,
		Picture:  picture,
		IsAdmin:  isAdmin,
	})
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        dbUser.ID,
		AuthID:    dbUser.AuthID,
		Email:     dbUser.Email,
		Name:      dbUser.Name,
		Picture:   dbUser.Picture.String,
		IsAdmin:   dbUser.IsAdmin.Bool,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}, nil
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	if r.queries == nil {
		return nil, nil // Database not connected
	}

	dbUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        dbUser.ID,
		AuthID:    dbUser.AuthID,
		Email:     dbUser.Email,
		Name:      dbUser.Name,
		Picture:   dbUser.Picture.String,
		IsAdmin:   dbUser.IsAdmin.Bool,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}, nil
}

// GetAllUsers retrieves all users
func (r *UserRepository) GetAllUsers(ctx context.Context) ([]User, error) {
	if r.queries == nil {
		return nil, nil // Database not connected
	}

	dbUsers, err := r.queries.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = User{
			ID:        dbUser.ID,
			AuthID:    dbUser.AuthID,
			Email:     dbUser.Email,
			Name:      dbUser.Name,
			Picture:   dbUser.Picture.String,
			IsAdmin:   dbUser.IsAdmin.Bool,
			CreatedAt: dbUser.CreatedAt.Time,
			UpdatedAt: dbUser.UpdatedAt.Time,
		}
	}

	return users, nil
}

// CountUsers returns the total number of users
func (r *UserRepository) CountUsers(ctx context.Context) (int64, error) {
	if r.queries == nil {
		return 0, nil // Database not connected
	}

	return r.queries.CountUsers(ctx)
}

// UpdateUser updates user information
func (r *UserRepository) UpdateUser(ctx context.Context, user *User) (*User, error) {
	if r.queries == nil {
		return nil, nil // Database not connected
	}

	picture := sql.NullString{String: user.Picture, Valid: user.Picture != ""}

	err := r.queries.UpdateUser(ctx, dbSqlc.UpdateUserParams{
		ID:       user.ID,
		Name:     user.Name,
		Picture:  picture,
	})
	if err != nil {
		return nil, err
	}

	// Return updated user by fetching it again
	return r.GetUserByEmail(ctx, user.Email)
}

// GetRecentUsers returns recently created users
func (r *UserRepository) GetRecentUsers(ctx context.Context) ([]User, error) {
	if r.queries == nil {
		return nil, nil // Database not connected
	}

	dbUsers, err := r.queries.GetRecentUsers(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = User{
			ID:        dbUser.ID,
			AuthID:    "", // Recent users query doesn't return auth_id
			Email:     dbUser.Email,
			Name:      dbUser.Name,
			Picture:   "", // Recent users query doesn't return picture
			IsAdmin:   false, // Recent users query doesn't return is_admin
			CreatedAt: dbUser.CreatedAt.Time,
			UpdatedAt: time.Time{}, // Recent users query doesn't return updated_at
		}
	}

	return users, nil
}

// CountUsersCreatedToday returns count of users created today
func (r *UserRepository) CountUsersCreatedToday(ctx context.Context) (int64, error) {
	if r.queries == nil {
		return 0, nil // Database not connected
	}

	return r.queries.CountUsersCreatedToday(ctx)
}

// CountUsersCreatedThisWeek returns count of users created this week
func (r *UserRepository) CountUsersCreatedThisWeek(ctx context.Context) (int64, error) {
	if r.queries == nil {
		return 0, nil // Database not connected
	}

	return r.queries.CountUsersCreatedThisWeek(ctx)
}
