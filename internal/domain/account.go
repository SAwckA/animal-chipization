package domain

type Account struct {
	ID        int32  `json:"id" db:"id"`
	FirstName string `json:"firstName" db:"firstname"`
	LastName  string `json:"lastName" db:"lastname"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"password" db:"password"`
}

type AccountResponse struct {
	ID        int32  `json:"id" db:"id"`
	FirstName string `json:"firstName" db:"firstname"`
	LastName  string `json:"lastName" db:"lastname"`
	Email     string `json:"email" db:"email"`
}

func MakeAccountResponse(account Account) *AccountResponse {
	return &AccountResponse{
		ID:        account.ID,
		FirstName: account.FirstName,
		LastName:  account.LastName,
		Email:     account.Email,
	}
}

func NewAccount(firstName, lastName, email, password string) *Account {
	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	}
}
