package usecase

import (
	"animal-chipization/internal/domain"
)

type accountRepository interface {
	GetByID(id int) (*domain.Account, error)
	Search(dto domain.SearchAccountDTO, size int, from int) (*[]domain.Account, error)
	Update(newAccount *domain.Account) error
	Delete(accoundID int) error
}

type AccountUsecase struct {
	repo accountRepository
}

func NewAccountUsecase(repo accountRepository) *AccountUsecase {
	return &AccountUsecase{
		repo: repo,
	}
}

// AccountUsecase.Get возвращает аккаунт, либо ошибки:
func (u *AccountUsecase) Get(id int) (*domain.Account, error) {
	return u.repo.GetByID(id)
}

// AccountUsecase.Search возращает список аккаунтов
// по критериям поиска domain.SearchAccountDTO или ошибки:
func (u *AccountUsecase) Search(dto domain.SearchAccountDTO, size int, from int) (*[]domain.Account, error) {

	if from < 0 {
		return nil, &domain.ApplicationError{
			OriginalError: nil,
			SimplifiedErr: domain.ErrInvalidInput,
			Description:   "invalid from param",
		}
	}
	if size <= 0 {
		return nil, &domain.ApplicationError{
			OriginalError: nil,
			SimplifiedErr: domain.ErrInvalidInput,
			Description:   "invalid size param",
		}
	}
	return u.repo.Search(dto, size, from)
}

// AccountUsecase.Update изменяет поля текущего аккаунта на новые
// и записывает в repository, возращает новый аккаунт и ошибки:
func (u *AccountUsecase) Update(currentAccount *domain.Account, newAccount domain.UpdateAccountDTO) (*domain.Account, error) {
	if currentAccount.ID != newAccount.ID {
		return nil, &domain.ApplicationError{
			OriginalError: nil,
			SimplifiedErr: domain.ErrForbidden,
			Description:   "update not your account",
		}
	}

	account := &domain.Account{
		ID:        newAccount.ID,
		FirstName: newAccount.FirstName,
		LastName:  newAccount.LastName,
		Email:     newAccount.Email,
		Password:  newAccount.Password,
	}

	return account, u.repo.Update(account)
}

// AccountUsecase.Delete удаляет аккаунт по id, возвращает ошибки:
func (u *AccountUsecase) Delete(executor *domain.Account, id int) error {
	if executor.ID != id {
		return &domain.ApplicationError{
			OriginalError: nil,
			SimplifiedErr: domain.ErrForbidden,
		}
	}
	return u.repo.Delete(id)
}
