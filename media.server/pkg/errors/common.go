package errors

import "net/http"

var (
	ErrInvalidCredentials = New(
		http.StatusUnauthorized,
		"INVALID_CREDENTIALS",
		"invalid credentials",
	)

	ErrUserAlreadyExists = New(
		http.StatusConflict,
		"USER_ALREADY_EXISTS",
		"user already exists",
	)

	ErrUnauthorized = New(
		http.StatusUnauthorized,
		"UNAUTHORIZED",
		"unauthorized",
	)

	ErrForbidden = New(
		http.StatusForbidden,
		"FORBIDDEN",
		"forbidden",
	)

	ErrNotFound = New(
		http.StatusNotFound,
		"NOT_FOUND",
		"resource not found",
	)
)
