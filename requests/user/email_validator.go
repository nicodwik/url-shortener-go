package requests

type UserEmailValidation struct {
	Email string `validate:"required,email"`
}
