package delivery

// NewUser ...
type NewUser struct {
	Username string `json:"username" binding:"required"`
}

// UserCreated ...
type UserCreated struct {
	ID string `json:"id"`
}
