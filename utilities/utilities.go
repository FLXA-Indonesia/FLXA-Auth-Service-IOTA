package utilities

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func ValidateJSON(err error) (string, int, []string) {
	message := ""
	code := 0
	var detail []string
	if err == nil {
		return message, code, detail
	}
	message = "Invalid JSON"
	isJsonEmpty := err.Error() == "EOF"
	if isJsonEmpty {
		detail = append(detail, "Body is empty")
		code = 400
		return message, code, detail
	}
	code = 500
	detail = append(detail, err.Error())
	return message, code, detail
}

func mapPgErrorToHTTPCode(pgCode string) int {
	switch pgCode {
	case "23505": // unique violation
		return 409
	case "23502": // not null violation
		return 400
	case "42601": // SQL syntax error
		return 500
	default:
		return 500
	}
}

func mapPgErrorToMessage(pgCode string) string {
	switch pgCode {
	case "23505": // unique violation
		return "Duplicate"
	case "23502": // not null violation
		return "Not Null Violation"
	case "42601": // SQL syntax error
		return "SQL Syntax Error"
	default:
		return "Transaction Error"
	}
}

func TransactionErrorHandler(transactionError error) (string, int, []string) {
	var pgErr *pgconn.PgError
	message := ""
	code := 0
	var details []string

	if transactionError == nil {
		return message, code, details
	}

	isPostgresError := errors.As(transactionError, &pgErr)
	if isPostgresError {
		message = mapPgErrorToMessage(pgErr.Code)
		details = append(details, pgErr.Message)
		code = mapPgErrorToHTTPCode(pgErr.Code)
	}
	return message, code, details
}

func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}
