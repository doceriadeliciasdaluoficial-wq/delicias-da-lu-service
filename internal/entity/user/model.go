package user

import "time"

type AdminUser struct {
	ID        string    `json:"id" firestore:"id"`
	Username  string    `json:"username" firestore:"username"`
	Email     string    `json:"email" firestore:"email"`
	Password  string    `json:"-" firestore:"password"`
	Role      string    `json:"role" firestore:"role"` // admin or manager
	LastLogin time.Time `json:"lastLogin" firestore:"lastLogin"`
	CreatedAt time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" firestore:"updatedAt"`
}

type LoginRequest struct {
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
}

type LoginResponse struct {
	Token string    `json:"token"`
	User  AdminUser `json:"user"`
}

type RefreshTokenResponse struct {
	Token string `json:"token"`
}
