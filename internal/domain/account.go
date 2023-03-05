package domain

const (
	AccountDefaultFrom = 0
	AccountDefaultSize = 10
)

type AccountID struct {
	ID int `uri:"accountId" binding:"gt=0"`
}

// var ErrAccountNotFoundByID = errors.New("account not found by id")                // Аккаунт с таким accountId не найден
// var ErrAccountLinkedWithAnimal = errors.New("account linked with animal")         // Аккаунт связан с животным
// var ErrAccountAccessForbidden = errors.New("delete non-personal account")         // Удаление не своего аккаунта
// var ErrAccountAlreadyExist = errors.New("account with given email already exist") // Аккаунт с таких email уже существует
// var ErrInvalidAccountForm = errors.New("invalid registration form")               // Невалидные данные для регистрации
// var ErrInvalidCredentials = errors.New("invalid credentials")

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

type SearchAccount struct {
	FirstName *string `form:"firstName"`
	LastName  *string `form:"lastName"`
	Email     *string `form:"email"`
	Size      *int    `form:"size"`
	From      *int    `form:"from"`
}

func (s *SearchAccount) Validate() error {
	err := &ApplicationError{
		OriginalError: nil,
		SimplifiedErr: ErrInvalidInput,
		Description:   "validation error",
	}
	var defaultFrom, defaultSize = AccountDefaultFrom, AccountDefaultSize

	if s.From == nil {
		s.From = &defaultFrom
	}
	if s.Size == nil {
		s.Size = &defaultSize
	}

	switch {
	case *s.From < 0:
		return err
	case *s.Size <= 0:
		return err

	default:
		return nil
	}
}

type UpdateAccount struct {
	ID        int
	FirstName string `json:"firstName" binding:"required,exclude_whitespace"`
	LastName  string `json:"lastName" binding:"required,exclude_whitespace"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,exclude_whitespace"`
}
