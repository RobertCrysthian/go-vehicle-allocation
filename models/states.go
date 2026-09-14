package models

type CreateStateDto struct {
	Id                int    `json:"id" validate:"required,number"`
	Name              string `json:"name" validate:"required"`
	StateAbbreviation string `json:"stateAbbreviation" validate:"required"`
}

type ListStateDto struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	StateAbbreviation string `json:"stateAbbreviation"`
}

type IbgeStateEntity struct {
	ID                int    `json:"id"`
	Name              string `json:"nome"`
	StateAbbreviation string `json:"sigla"`
}
