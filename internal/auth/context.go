package auth

type ContextKey string

const (
	ContextKeyRole      ContextKey = "role"
	ContextKeyExpiresAt ContextKey = "expiresAt"
	ContextKeyUsername  ContextKey = "username"
	ContextKeyUserID    ContextKey = "userID"
)
