package usecase

import (
	"animal-chipization/internal/domain"

	"github.com/go-playground/validator/v10"
)

type registerAccountReporitory interface {
	Create(account *domain.Account) (int, error)
	GetByEmail(email string) (*domain.Account, error)
}

type RegisterAccountUsecase struct {
	repo registerAccountReporitory
}

func NewRegisterAccountUsecase(repo registerAccountReporitory) *RegisterAccountUsecase {
	return &RegisterAccountUsecase{
		repo: repo,
	}
}

// RegisterAccountUsecase.Register регистрирует новый аккаунт,
// возвращает новосозданный аккаунт и ошибки:
//		domain.ErrAccountAlreadyExist :: repository
// 		domain.ErrInvalidAccountForm  :: validation
func (u *RegisterAccountUsecase) Register(dto domain.RegistrationDTO) (*domain.Account, error) {

	if err := validator.New().Struct(dto); err != nil {
		return nil, err
	}

	account := domain.NewAccount(dto)

	id, err := u.repo.Create(account)
	account.ID = id

	return account, err
}

// RegisterAccountUsecase.Login аутентификация пользователя,
// возращает аккаунт под которым авторизовались и ошибки:
// 		domain.ErrAccountNotFoundByID :: repository
//
// 		domain.ErrInvalidCredentials  :: usecase
func (u *RegisterAccountUsecase) Login(email, password string) (*domain.Account, error) {
	account, err := u.repo.GetByEmail(email)

	if err != nil {
		return nil, err
	}

	if account.Password == password {
		return account, nil
	}

	return nil, domain.ErrInvalidCredentials
}
