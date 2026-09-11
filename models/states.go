package models

type CreateStateDto struct {
	Name              string `json:"name" validate:"required"`
	StateAbbreviation string `json:"stateAbbreviation" validate:"required"`
	IbgeId            int    `json:"ibgeId" validate:"required,number"`
}

type ListStateDto struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	StateAbbreviation string `json:"stateAbbreviation"`
	IbgeId            int    `json:"ibgeId"`
}
