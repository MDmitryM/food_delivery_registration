package models

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UpdateUser struct {
	ID      int32  `json:"id"`
	PwdHash string `json:"pwd_hash"`
}
