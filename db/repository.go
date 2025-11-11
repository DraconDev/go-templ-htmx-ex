package db

import (
	"database/sql"
	"fmt"
)

// UserRepository provides database operations for users
type UserRepository struct {
	db *Database
}

// NewUserRepository creates a new user repository
func NewUserRepository(database *Database) *UserRepository {
	return &UserRepository{db: database}
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(user *User) error {
	query := `
		INSERT INTO users (auth_id, email, name, picture)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(query, user.AuthID, user.Email, user.Name, user.Picture).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByAuthID retrieves a user by their auth service ID
func (r *UserRepository) GetUserByAuthID(authID string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, auth_id, email, name, picture, created_at, updated_at
		FROM users
		WHERE auth_id = $1
	`

	err := r.db.QueryRow(query, authID).
		Scan(&user.ID, &user.AuthID, &user.Email, &user.Name, &user.Picture, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil // User not found
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user by auth_id: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by their email
func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, auth_id, email, name, picture, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := r.db.QueryRow(query, email).
		Scan(&user.ID, &user.AuthID, &user.Email, &user.Name, &user.Picture, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil // User not found
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by their ID
func (r *UserRepository) GetUserByID(id string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, auth_id, email, name, picture, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).
		Scan(&user.ID, &user.AuthID, &user.Email, &user.Name, &user.Picture, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil // User not found
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

// UpdateUser updates user information
func (r *UserRepository) UpdateUser(user *User) error {
	query := `
		UPDATE users
		SET email = $2, name = $3, picture = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRow(query, user.ID, user.Email, user.Name, user.Picture).
		Scan(&user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// GetAllUsers retrieves all users for admin dashboard
func (r *UserRepository) GetAllUsers() ([]User, error) {
	users := []User{}
	query := `
		SELECT id, auth_id, email, name, picture, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.AuthID, &user.Email, &user.Name, &user.Picture, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUserCount returns the total number of users
func (r *UserRepository) GetUserCount() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users`

	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get user count: %w", err)
	}

	return count, nil
}

// GetNewUsersCount returns count of users created in the last N days
func (r *UserRepository) GetNewUsersCount(days int) (int, error) {
	var count int
	query := `
		SELECT COUNT(*) FROM users
		WHERE created_at >= NOW() - INTERVAL '1 day' * $1
	`

	err := r.db.QueryRow(query, days).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get new users count: %w", err)
	}

	return count, nil
}

// UserSession represents a user session for admin analytics
type UserSession struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	LoginTime time.Time `json:"login_time"`
	IPAddress string    `json:"ip_address"`
	Active    bool      `json:"active"`
}

// GetActiveSessions returns currently active user sessions
func (r *UserRepository) GetActiveSessions() ([]UserSession, error) {
	sessions := []UserSession{}
	query := `
		SELECT s.id, s.user_id, u.email, u.name, s.login_time, s.ip_address, s.active
		FROM user_sessions s
		JOIN users u ON s.user_id = u.id
		WHERE s.active = true AND s.last_seen > NOW() - INTERVAL '1 hour'
		ORDER BY s.last_seen DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active sessions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var session UserSession
		err := rows.Scan(&session.ID, &session.UserID, &session.Email, &session.Name, &session.LoginTime, &session.IPAddress, &session.Active)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// UserPreferencesRepository provides database operations for user preferences
type UserPreferencesRepository struct {
	db *Database
}

// NewUserPreferencesRepository creates a new user preferences repository
func NewUserPreferencesRepository(database *Database) *UserPreferencesRepository {
	return &UserPreferencesRepository{db: database}
}

// CreateUserPreferences creates default user preferences
func (r *UserPreferencesRepository) CreateUserPreferences(prefs *UserPreferences) error {
	query := `
		INSERT INTO user_preferences (user_id, theme, language, email_notifications, push_notifications)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(query, prefs.UserID, prefs.Theme, prefs.Language, prefs.EmailNotifications, prefs.PushNotifications).
		Scan(&prefs.ID, &prefs.CreatedAt, &prefs.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user preferences: %w", err)
	}

	return nil
}

// GetUserPreferences retrieves user preferences
func (r *UserPreferencesRepository) GetUserPreferences(userID string) (*UserPreferences, error) {
	prefs := &UserPreferences{}
	query := `
		SELECT id, user_id, theme, language, email_notifications, push_notifications, created_at, updated_at
		FROM user_preferences
		WHERE user_id = $1
	`

	err := r.db.QueryRow(query, userID).
		Scan(&prefs.ID, &prefs.UserID, &prefs.Theme, &prefs.Language, &prefs.EmailNotifications, &prefs.PushNotifications, &prefs.CreatedAt, &prefs.UpdatedAt)

	if err == sql.ErrNoRows {
		// Create default preferences if not found
		defaultPrefs := &UserPreferences{
			UserID:             userID,
			Theme:              "dark",
			Language:           "en",
			EmailNotifications: true,
			PushNotifications:  true,
		}

		err := r.CreateUserPreferences(defaultPrefs)
		if err != nil {
			return nil, fmt.Errorf("failed to create default user preferences: %w", err)
		}

		return defaultPrefs, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user preferences: %w", err)
	}

	return prefs, nil
}

// UpdateUserPreferences updates user preferences
func (r *UserPreferencesRepository) UpdateUserPreferences(prefs *UserPreferences) error {
	query := `
		UPDATE user_preferences
		SET theme = $2, language = $3, email_notifications = $4, push_notifications = $5, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRow(query, prefs.ID, prefs.Theme, prefs.Language, prefs.EmailNotifications, prefs.PushNotifications).
		Scan(&prefs.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update user preferences: %w", err)
	}

	return nil
}
