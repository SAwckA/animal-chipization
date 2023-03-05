package domain

import "errors"

var ErrInvalidInput = errors.New("invalid input")
var ErrLinked = errors.New("entity linked")
var ErrConflict = errors.New("data conflict")
var ErrAlreadyExist = errors.New("already exist")
var ErrBadDatabaseOut = errors.New("invalid database out")
var ErrNotFound = errors.New("not found")
var ErrUnknown = errors.New("unknown error")
var ErrForbidden = errors.New("forbidden")

// ApplicationError обогощение ошибки, для упрощенной обработки в контроллерах
type ApplicationError struct {
	// Ошибка не связанная с логикой функции
	// (Сторонняя ошибка)
	OriginalError error

	// Ошибка читаемая приложением
	// Она же упрощенная
	SimplifiedErr error

	// Подробное описание экземпляра ошибки
	Description string
}

func (e *ApplicationError) Error() string {
	if e.OriginalError == nil && e.SimplifiedErr == nil {
		return e.Description
	}

	if e.OriginalError == nil {
		return e.SimplifiedErr.Error() + ": " + e.Description
	}

	if e.SimplifiedErr == nil {
		return e.OriginalError.Error() + ": " + e.Description
	}

	return e.OriginalError.Error() + " -> " + e.SimplifiedErr.Error() + ": " + e.Description
}

func (e *ApplicationError) Unwrap() error {
	return e.SimplifiedErr
}
