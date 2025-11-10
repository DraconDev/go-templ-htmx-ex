package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Profile represents application-specific user profile data
// Complements the external auth service by storing app-level data
type Profile struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"` // References auth service user ID
	Bio         string    `json:"bio" db:"bio"`
	Website     string    `json:"website" db:"website"`
	Company     string    `json:"company" db:"company"`
	Location    string    `json:"location" db:"location"`
	Timezone    string    `json:"timezone" db:"timezone"`
	Language    string    `json:"language" db:"language"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// UserPreferences represents user application preferences
// Global settings that persist across sessions
type UserPreferences struct {
	ID               string `json:"id" db:"id"`
	UserID           string `json:"user_id" db:"user_id"`
	Theme            string `json:"theme" db:"theme"` // "dark", "light", "auto"
	Language         string `json:"language" db:"language"` // "en", "es", "fr", etc.
	EmailNotifications bool `json:"email_notifications" db:"email_notifications"`
	PushNotifications bool `json:"push_notifications" db:"push_notifications"`
	DashboardLayout  string `json:"dashboard_layout" db:"dashboard_layout"` // JSON layout config
	TimeFormat       string `json:"time_format" db:"time_format"` // "12h", "24h"
	DateFormat       string `json:"date_format" db:"date_format"` // "US", "EU", "ISO"
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// Session represents application-level user sessions
// Complements JWT auth with app-specific session tracking
type Session struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	JWTPayload  string    `json:"jwt_payload" db:"jwt_payload"` // Store JWT claims for reference
	DeviceInfo  string    `json:"device_info" db:"device_info"` // JSON device/browser info
	IPAddress   string    `json:"ip_address" db:"ip_address"`
	Location    string    `json:"location" db:"location"`
	LastActivity time.Time `json:"last_activity" db:"last_activity"`
	ExpiresAt   time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	IsActive    bool      `json:"is_active" db:"is_active"`
}

// AuditLog represents user and system actions for compliance and debugging
type AuditLog struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	Action      string    `json:"action" db:"action"` // "login", "logout", "update_profile", etc.
	Resource    string    `json:"resource" db:"resource"` // "user", "profile", "settings", etc.
	ResourceID  string    `json:"resource_id" db:"resource_id"`
	OldValues   JSONMap   `json:"old_values" db:"old_values"` // Previous values (for updates)
	NewValues   JSONMap   `json:"new_values" db:"new_values"` // New values (for updates)
	IPAddress   string    `json:"ip_address" db:"ip_address"`
	UserAgent   string    `json:"user_agent" db:"user_agent"`
	Success     bool      `json:"success" db:"success"` // Whether action succeeded
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// ApplicationMetadata stores app-level configuration and metadata
type ApplicationMetadata struct {
	Key         string    `json:"key" db:"key"`
	Value       JSONMap   `json:"value" db:"value"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// JSONMap is a custom type for JSON data storage
type JSONMap map[string]interface{}

// Scan implements the Scanner interface for JSONMap
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value implements the Valuer interface for JSONMap
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Database schema - designed for startup flexibility
const DatabaseSchema = `
-- Profiles table: Application-specific user data
-- Complements external auth service
CREATE TABLE IF NOT EXISTS profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) UNIQUE NOT NULL, -- References auth service user ID
    bio TEXT,
    website VARCHAR(500),
    company VARCHAR(255),
    location VARCHAR(255),
    timezone VARCHAR(50) DEFAULT 'UTC',
    language VARCHAR(10) DEFAULT 'en',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- User preferences: Global user settings
CREATE TABLE IF NOT EXISTS user_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) UNIQUE NOT NULL,
    theme VARCHAR(20) DEFAULT 'dark',
    language VARCHAR(10) DEFAULT 'en',
    email_notifications BOOLEAN DEFAULT true,
    push_notifications BOOLEAN DEFAULT true,
    dashboard_layout JSONB DEFAULT '{}',
    time_format VARCHAR(5) DEFAULT '24h',
    date_format VARCHAR(10) DEFAULT 'ISO',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Sessions: Application-level session tracking
-- Complements JWT with app-specific session management
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    jwt_payload TEXT, -- Store JWT claims for reference
    device_info JSONB DEFAULT '{}',
    ip_address INET,
    location VARCHAR(255),
    last_activity TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_active BOOLEAN DEFAULT true
);

-- Audit logs: Track user actions for compliance and debugging
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255),
    action VARCHAR(100) NOT NULL,
    resource VARCHAR(100) NOT NULL,
    resource_id VARCHAR(255),
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    success BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Application metadata: Store app configuration
CREATE TABLE IF NOT EXISTS application_metadata (
    key VARCHAR(255) PRIMARY KEY,
    value JSONB NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_profiles_user_id ON profiles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_preferences_user_id ON user_preferences(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_last_activity ON sessions(last_activity);
CREATE INDEX IF NOT EXISTS idx_sessions_is_active ON sessions(is_active);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for updated_at
CREATE TRIGGER update_profiles_updated_at 
    BEFORE UPDATE ON profiles
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_preferences_updated_at 
    BEFORE UPDATE ON user_preferences
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_application_metadata_updated_at 
    BEFORE UPDATE ON application_metadata
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default application metadata
INSERT INTO application_metadata (key, value, description) VALUES 
('app_name', '"Startup Platform"', 'Application name'),
('app_version', '"1.0.0"', 'Application version'),
('features', '{"auth": true, "profiles": true, "preferences": true, "audit": true}', 'Enabled features'),
('limits', '{"max_sessions": 10, "max_file_size": "10MB", "rate_limit": 100}', 'Application limits')
ON CONFLICT (key) DO NOTHING;
`