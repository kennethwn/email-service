package error_wrap

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrSqlError            = errors.New("database server failed to execute query")
	ErrTooManyRequests     = errors.New("too many requests")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInvalidToken        = errors.New("invalid token")
	ErrApiKeyIsMissing     = errors.New("api key is missing")
	ErrApiKeyIsInvalid     = errors.New("api key is invalid")
	ErrForbidden           = errors.New("forbidden")
	ErrBadRequest          = errors.New("bad request")
	ErrNotFound            = errors.New("not found")
	ErrIPorServiceBlocked  = errors.New("ip or service is blocked")
)

var GeneralErrors = []error{
	ErrInternalServerError,
	ErrSqlError,
	ErrTooManyRequests,
	ErrUnauthorized,
	ErrInvalidToken,
	ErrApiKeyIsMissing,
	ErrApiKeyIsInvalid,
	ErrForbidden,
	ErrBadRequest,
	ErrNotFound,
	ErrIPorServiceBlocked,
}
