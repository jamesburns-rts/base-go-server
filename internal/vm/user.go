package vm

type (
	Login struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=72"`
	}
	Tokens struct {
		AccessToken string `json:"accessToken"`
	}
)

type Example struct {
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Moves []string `json:"moves"`
}
