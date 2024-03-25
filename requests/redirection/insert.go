package requests

type InsertRedirectionValidation struct {
	// ShortUrl string `validate:"required"`
	LongUrl string `validate:"required"`
}
