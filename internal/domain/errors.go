package domain

import "errors"

var ErrInvalidInput = errors.New("invalid input") // Некорректные входные данных (нельзя обработать, ввиду логики)
var ErrConflict = errors.New("data conflict")     // Конфликт данных
var ErrAlreadyExist = errors.New("already exist") // Уже существует
var ErrNotFound = errors.New("not found")         // Сущность не найдена
var ErrForbidden = errors.New("forbidden")        // Доступ запрещён
var ErrUnknown = errors.New("unknown error")      // Неизвестная ошибка

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

	return "[" + e.SimplifiedErr.Error() + "] " + e.Description
}

func (e *ApplicationError) Unwrap() error {
	return e.SimplifiedErr
}
