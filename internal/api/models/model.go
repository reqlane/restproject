package models

type ModelWithID interface {
	GetID() int
}

type ModelWithPassword interface {
	GetPassword() string
	SetPassword(password string)
}
