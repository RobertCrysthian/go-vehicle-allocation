package handlers

import (
	"encoding/json"
	"net/http"
)

type AboutResponse struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Features      []string `json:"features"`
	BusinessRules []string `json:"businessRules"`
	InDevelopment string   `json:"inDevelopment"`
}

func About(writer http.ResponseWriter, request *http.Request) {
	response := AboutResponse{
		Name: "go-vehicle-allocation",
		Description: "Sistema para alocação de carros. Permite criar um usuário, selecionar um " +
			"intervalo de data, selecionar um carro disponível e alocá-lo (criando um registro de " +
			"alocação) para o intervalo selecionado.",
		Features: []string{
			"Cadastro de usuários, com validação de duplicidade de CPF e email",
			"Rota simples de login capaz de validar os dados do usuário (autenticação ainda em progresso)",
			"Seeders de cidades e estados, populados a partir da api de localidades do IBGE",
			"Pontos de coleta (locais onde o cliente busca o veículo alocado), criados a partir das cidades e estados",
			"Existe uma modelagem criada através do draw.io para representar as tabelas e seus relacionamentos",
			"A query de criação está dentro de seeders",
		},
		BusinessRules: []string{
			"Um veículo deve ter um ponto de coleta, e um ponto de coleta pode estar em diversos veículos",
			"Cada veículo possui um valor de alocação, cobrado por dia alocado, e uma taxa de atraso, cobrada pelos dias que excedem a data final estimada da alocação",
			"Ao selecionar veículos por intervalo de datas, devem ser exibidos apenas os que não possuem alocações em aberto naquele período",
			"Ao listar as alocações, deve ser exibido o valor total (dias * preço do veículo) e o valor da taxa de atraso (taxa * dias atrasados)",
		},
		InDevelopment: "Ainda há bastante coisa para desenvolver para isso funcionar por completo; o progresso continua até o prazo expirar.",
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}
