package error

import "errors"

var (
	ErrNoFields          = errors.New("no author or text fields required")
	ErrAuthorNotFound    = errors.New("author not found")
	ErrRecordNotFound    = errors.New("record not found")
	ErrNoQuotesAvailable = errors.New("no records")
	ErrInvalidTypeID     = errors.New("id can only be a number")

	ErrTypeConversion = errors.New("type conversion error")
)
