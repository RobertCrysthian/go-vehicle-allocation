package models

type CreateUserDto struct {
	Name     string `json:"name" validate:"required,min=4,max=32"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
	CPF      string `json:"cpf" validate:"len=11,number"`
}

// type UpdateUser struct {
// 	CreateUser
// }

type ListUsersDto struct {
	ID    int    `json:"id"`
	Name  string `json:"title"`
	Email string `json:"email"`
	CPF   string `json:"cpf"`
}