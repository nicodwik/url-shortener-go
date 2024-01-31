package requests

type InsertUrlValidation struct {
	ShortUrl string `validate:"required"`
	LongUrl  string `validate:"required"`
}
