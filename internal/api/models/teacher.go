package models

type Teacher struct {
	ID        int    `json:"id,omitempty" db:"id"`
	FirstName string `json:"first_name,omitempty" db:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name,omitempty" db:"last_name" validate:"required,min=2,max=50"`
	Email     string `json:"email,omitempty" db:"email" validate:"required,email"`
	Class     string `json:"class,omitempty" db:"class" validate:"required"`
	Subject   string `json:"subject,omitempty" db:"subject" validate:"required"`
}

func (t *Teacher) GetID() int {
	return t.ID
}
