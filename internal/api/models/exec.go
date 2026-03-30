package models

import "database/sql"

type Exec struct {
	ID                   int            `json:"id,omitempty" db:"id"`
	FirstName            string         `json:"first_name,omitempty" db:"first_name"`
	LastName             string         `json:"last_name,omitempty" db:"last_name"`
	Email                string         `json:"email,omitempty" db:"email"`
	Username             string         `json:"username,omitempty" db:"username"`
	Password             string         `json:"password,omitempty" db:"password"`
	PasswordChangedAt    sql.NullString `json:"password_changed_at,omitzero" db:"password_changed_at"`
	UserCreatedAt        sql.NullString `json:"user_created_at,omitzero" db:"user_created_at,omitinsert"`
	PasswordResetToken   sql.NullString `json:"password_reset_token,omitzero" db:"password_reset_token"`
	PasswordTokenExpires sql.NullString `json:"password_token_expires,omitzero" db:"password_token_expires"`
	InactiveStatus       bool           `json:"inactive_status" db:"inactive_status"`
	Role                 string         `json:"role,omitempty" db:"role"`
}

func (e *Exec) GetID() int {
	return e.ID
}

type ExecResponse struct {
	ID             int    `json:"id,omitempty"`
	FirstName      string `json:"first_name,omitempty"`
	LastName       string `json:"last_name,omitempty"`
	Email          string `json:"email,omitempty"`
	Username       string `json:"username,omitempty"`
	UserCreatedAt  string `json:"user_created_at,omitempty"`
	InactiveStatus bool   `json:"inactive_status"`
	Role           string `json:"role,omitempty"`
}

func (e *Exec) ToResponse() *ExecResponse {
	return &ExecResponse{
		ID:             e.ID,
		FirstName:      e.FirstName,
		LastName:       e.LastName,
		Email:          e.Email,
		Username:       e.Username,
		UserCreatedAt:  e.UserCreatedAt.String,
		InactiveStatus: e.InactiveStatus,
		Role:           e.Role,
	}
}

type Execs []Exec

func (execs Execs) ToResponse() []ExecResponse {
	responses := make([]ExecResponse, len(execs))
	for i := range execs {
		responses[i] = *execs[i].ToResponse()
	}
	return responses
}

type ExecCredentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type UpdatePasswordResponse struct {
	Token           string `json:"token"`
	PasswordUpdated bool   `json:"password_updated"`
}
