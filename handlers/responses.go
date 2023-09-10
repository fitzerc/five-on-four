package handlers

import "github.com/fitzerc/five-on-four/data"

type LoginResponse struct {
	Email     string          `json:"email"`
	Password  string          `json:"password"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	Token     string          `json:"token"`
	Picture   []byte          `json:"picture"`
	Roles     []data.UserRole `json:"roles"`
}
