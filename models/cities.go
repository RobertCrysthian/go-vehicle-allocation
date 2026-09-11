package models

type CreateCityDto struct {
	Name string `json:"name" validate:"required"`

	IbgeId int `json:"ibgeId" validate:"required,number"`
}

type ListCityDto struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	IbgeId int    `json:"ibgeId"`
}
