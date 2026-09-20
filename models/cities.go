package models

type CreateCityDto struct {
	Id      int    `json:"id" validate:"required,number"`
	Name    string `json:"name" validate:"required"`
	StateId int    `json:"stateId" validate:"required,number"`
}

type ListCityDto struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	StateAbbreviation string `json:"StateAbbreviation"`
}

type ListCityWithPaginationDto struct {
	Page    int           `json:"page"`
	Results int           `json:"results"`
	Total   int           `json:"total"`
	Cities  []ListCityDto `json:"cities"`
}

type ibgeUFEntity struct {
	ID int `json:"id"`
}

type IbgeCityEntity struct {
	ID           int    `json:"id"`
	Name         string `json:"nome"`
	Microrregiao struct {
		Mesorregiao struct {
			UF ibgeUFEntity `json:"UF"`
		} `json:"mesorregiao"`
	} `json:"microrregiao"`
	RegiaoImediata struct {
		RegiaoIntermediaria struct {
			UF ibgeUFEntity `json:"UF"`
		} `json:"regiao-intermediaria"`
	} `json:"regiao-imediata"`
}