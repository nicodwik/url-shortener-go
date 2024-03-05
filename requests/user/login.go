package requests

type LoginValidation struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}
