package models

import "time"

// UserPreferences represents user settings
type UserPreferences struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	Theme              string    `json:"theme"`
	Language           string    `json:"language"`
	Timezone           string    `json:"timezone"`
	EmailNotifications bool      `json:"email_notifications"`
	EmailBilling       bool      `json:"email_billing"`
	PushNotifications  bool      `json:"push_notifications"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
