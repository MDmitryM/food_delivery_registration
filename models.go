package models

type User struct {
	Login    string `json:"login" validate:"required" example:"login"`
	Password string `json:"password" validate:"required" example:"password"`
}

type UpdateUser struct {
	ID       int32  `json:"user_id"`
	Password string `json:"password" validate:"required"`
}
