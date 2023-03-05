package usecase

import (
	"animal-chipization/internal/domain"
)

type accountRepository interface {
	GetByID(id int) (*domain.Account, error)
	Search(dto *domain.SearchAccount) ([]domain.Account, error)
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

// Get возвращает аккаунт, либо ошибки:
func (u *AccountUsecase) Get(id int) (*domain.Account, error) {
	return u.repo.GetByID(id)
}

// Search возращает список аккаунтов
// по критериям поиска domain.SearchAccountDTO или ошибки:
func (u *AccountUsecase) Search(dto *domain.SearchAccount) ([]domain.Account, error) {
	if err := dto.Validate(); err != nil {
		return nil, err
	}

	return u.repo.Search(dto)
}

// Update изменяет поля текущего аккаунта на новые
// и записывает в repository, возращает новый аккаунт и ошибки:
func (u *AccountUsecase) Update(currentAccount *domain.Account, newAccount domain.UpdateAccount) (*domain.Account, error) {
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

// Delete удаляет аккаунт по id, возвращает ошибки:
func (u *AccountUsecase) Delete(executor *domain.Account, id int) error {
	if executor.ID != id {
		return &domain.ApplicationError{
			OriginalError: nil,
			SimplifiedErr: domain.ErrForbidden,
		}
	}
	return u.repo.Delete(id)
}
