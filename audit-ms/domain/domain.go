package domain

import "time"

type AuditLog struct {
	ID            string    `json:"id" db:"id"`
	UserType      string    `json:"user_type" db:"user_type"`
	UserID        string    `json:"user_id" db:"user_id"`
	Event         string    `json:"event" db:"event"`
	AuditableType string    `json:"auditable_type" db:"auditable_type"`
	AuditableID   string    `json:"auditable_id" db:"auditable_id"`
	OldValues     string    `json:"old_values" db:"old_values"`
	NewValues     string    `json:"new_values" db:"new_values"`
	URL           string    `json:"url" db:"url"`
	IPAddress     string    `json:"ip_address" db:"ip_address"`
	UserAgent     string    `json:"user_agent" db:"user_agent"`
	Tags          string    `json:"tags" db:"tags"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
