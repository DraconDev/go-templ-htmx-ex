package repositories

import (
	"context"

	dbSqlc "github.com/DraconDev/go-templ-htmx-ex/db/sqlc"
	"github.com/DraconDev/go-templ-htmx-ex/models"
)

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
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if r.queries == nil {
		return nil, models.ErrDatabaseNotConnected
	}

	dbUser, err := r.queries.CreateUser(ctx, dbSqlc.CreateUserParams{
		Email:    user.Email,
		Name:     user.Name,
		Picture:  user.Picture,
		IsAdmin:  user.IsAdmin,
		Provider: user.Provider,
	})
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		Name:      dbUser.Name,
		Picture:   dbUser.Picture.String,
		IsAdmin:   dbUser.IsAdmin.Bool,
		Provider:  dbUser.Provider.String,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}, nil
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if r.queries == nil {
		return nil, models.ErrDatabaseNotConnected
	}

	dbUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		Name:      dbUser.Name,
		Picture:   dbUser.Picture.String,
		IsAdmin:   dbUser.IsAdmin.Bool,
		Provider:  dbUser.Provider.String,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}, nil
}

// GetAllUsers retrieves all users
func (r *UserRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	if r.queries == nil {
		return nil, models.ErrDatabaseNotConnected
	}

	dbUsers, err := r.queries.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]models.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = models.User{
			ID:        dbUser.ID,
			Email:     dbUser.Email,
			Name:      dbUser.Name,
			Picture:   dbUser.Picture.String,
			IsAdmin:   dbUser.IsAdmin.Bool,
			Provider:  dbUser.Provider.String,
			CreatedAt: dbUser.CreatedAt.Time,
			UpdatedAt: dbUser.UpdatedAt.Time,
		}
	}

	return users, nil
}

// CountUsers returns the total number of users
func (r *UserRepository) CountUsers(ctx context.Context) (int64, error) {
	if r.queries == nil {
		return 0, models.ErrDatabaseNotConnected
	}

	return r.queries.CountUsers(ctx)
}

// UpdateUser updates user information
func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if r.queries == nil {
		return nil, models.ErrDatabaseNotConnected
	}

	dbUser, err := r.queries.UpdateUser(ctx, dbSqlc.UpdateUserParams{
		ID:       user.ID,
		Name:     user.Name,
		Picture:  user.Picture,
		IsAdmin:  user.IsAdmin,
		Provider: user.Provider,
	})
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		Name:      dbUser.Name,
		Picture:   dbUser.Picture.String,
		IsAdmin:   dbUser.IsAdmin.Bool,
		Provider:  dbUser.Provider.String,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}, nil
}

// GetRecentUsers returns recently created users
func (r *UserRepository) GetRecentUsers(ctx context.Context) ([]models.User, error) {
	if r.queries == nil {
		return nil, models.ErrDatabaseNotConnected
	}

	dbUsers, err := r.queries.GetRecentUsers(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]models.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = models.User{
			ID:        dbUser.ID,
			Email:     dbUser.Email,
			Name:      dbUser.Name,
			Picture:   dbUser.Picture.String,
			IsAdmin:   dbUser.IsAdmin.Bool,
			Provider:  dbUser.Provider.String,
			CreatedAt: dbUser.CreatedAt.Time,
			UpdatedAt: dbUser.UpdatedAt.Time,
		}
	}

	return users, nil
}

// CountUsersCreatedToday returns count of users created today
func (r *UserRepository) CountUsersCreatedToday(ctx context.Context) (int64, error) {
	if r.queries == nil {
		return 0, models.ErrDatabaseNotConnected
	}

	return r.queries.CountUsersCreatedToday(ctx)
}

// CountUsersCreatedThisWeek returns count of users created this week
func (r *UserRepository) CountUsersCreatedThisWeek(ctx context.Context) (int64, error) {
	if r.queries == nil {
		return 0, models.ErrDatabaseNotConnected
	}

	return r.queries.CountUsersCreatedThisWeek(ctx)
}
