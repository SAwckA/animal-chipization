package domain

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

const (
	AccountDefaultFrom = 0
	AccountDefaultSize = 10
)

type AccountID struct {
	ID int `uri:"accountId" binding:"gt=0"`
}

var ErrAccountNotFoundByID = errors.New("account not found by id")                // Аккаунт с таким accountId не найден
var ErrAccountLinkedWithAnimal = errors.New("account linked with animal")         // Аккаунт связан с животным
var ErrAccountAccessForbidden = errors.New("delete non-personal account")         // Удаление не своего аккаунта
var ErrAccountAlreadyExist = errors.New("account with given email already exist") // Аккаунт с таких email уже существует
var ErrInvalidAccountForm = errors.New("invalid registration form")               // Невалидные данные для регистрации
var ErrInvalidCredentials = errors.New("invalid credentials")

type Account struct {
	ID        int    `json:"id" db:"id"`
	FirstName string `json:"firstName" db:"firstname"`
	LastName  string `json:"lastName" db:"lastname"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"password" db:"password"`
}

func (a *Account) Response() map[string]interface{} {
	return map[string]interface{}{
		"id":        a.ID,
		"firstName": a.FirstName,
		"lastName":  a.LastName,
		"email":     a.Email,
	}
}

type RegistrationDTO struct {
	FirstName string `json:"firstName" binding:"required,exclude_whitespace"`
	LastName  string `json:"lastName" binding:"required,exclude_whitespace"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,exclude_whitespace"`
}

func NewAccount(dto RegistrationDTO) *Account {
	return &Account{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Password:  dto.Password,
	}
}

type SearchAccountDTO struct {
	FirstName *string `form:"firstName"`
	LastName  *string `form:"lastName"`
	Email     *string `form:"email"`
}

func (d *SearchAccountDTO) Prepare() error {
	if err := validator.New().Struct(d); err != nil {
		return err
	}

	if d.FirstName != nil || *d.FirstName == "null" {
		d.FirstName = nil
	}

	if d.LastName != nil || *d.LastName == "null" {
		d.LastName = nil
	}

	if d.Email != nil || *d.Email == "null" {
		d.Email = nil
	}

	return nil
}

type UpdateAccountDTO struct {
	ID        int
	FirstName string `json:"firstName" binding:"required,exclude_whitespace"`
	LastName  string `json:"lastName" binding:"required,exclude_whitespace"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,exclude_whitespace"`
}
