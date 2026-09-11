package models

type CreateSessionDto struct {
	UserId    int    `json:"userId" validate:"required,number"`
	LoginAt   string `json:"loginAt" validate:"required,date"`
	ExpiredAt string `json:"expiredAt" validate:"required,date"`
	LogoutAt  string `json:"logoutAt" validate:"required,date"`
}
