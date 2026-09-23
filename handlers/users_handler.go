package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/robertcrysthian/backend-go/models"
	"github.com/robertcrysthian/backend-go/utils"
	"golang.org/x/crypto/bcrypt"
)

type UsersHandler struct {
	DB *sql.DB
}

func NewUsersHandler(db *sql.DB) *UsersHandler {
	return &UsersHandler{DB: db}
}

func (UsersHandler *UsersHandler) CreateUser (writer http.ResponseWriter, request *http.Request) {
	var user models.CreateUserDto
	err := utils.ValidateRequest(request, &user)
	if err != nil {
		utils.BadRequestError(writer, err.Error())
		return
	}

	const checkCpfAndEmailAvailabilityQuery = `
		SELECT EXISTS (
			SELECT 1 FROM users
			WHERE cpf = $1 OR email = $2
		)
	`

	var existingUserId bool

	if err := UsersHandler.DB.QueryRow(checkCpfAndEmailAvailabilityQuery, user.CPF, user.Email).Scan(&existingUserId); err != nil {
		utils.InternalServerError(writer, "Erro na query que valida email e cpf " +err.Error())
		return
	}

	if(existingUserId) {
		utils.ConflitctError(writer, "Esse usuário já está cadastrado no sistema!")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalServerError(writer, "Falha ao gerar hash da senha " +err.Error())
		return
	}

	const createAccountQuery = `
		INSERT INTO users (name, email, password, cpf) VALUES ($1, $2, $3, $4)
	`

	_, err = UsersHandler.DB.Exec(createAccountQuery, user.Name, user.Email, passwordHash, user.CPF)
	if err != nil {
		utils.InternalServerError(writer, "Erro na query de inserir usuário " +err.Error())
		return
	}

	response := map[string]string{"response": "Usuário cadastrado com sucesso!"}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}

//TODO - Precisa de uma autenticação REAL no projeto. Estarei trabalhando nisso
func (UsersHandler *UsersHandler) Login (writer http.ResponseWriter, request *http.Request) {
	var credentials models.LoginDto
	err := utils.ValidateRequest(request, &credentials)
	if err != nil {
		utils.BadRequestError(writer, err.Error())
		return
	}

	const query = `SELECT * from users where email = $1`
	var user models.UserWithHashDto
	err = UsersHandler.DB.QueryRow(query, credentials.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CPF)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.BadRequestError(writer, "Usuário ou senha inválidos")
			return
		}
		utils.InternalServerError(writer, "Ocorreu um erro ao escanear o veículo: " + err.Error())
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		utils.UnauthorizedError(writer, "Usuário ou senha inválidos")
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(user.ListUsersDto)
}