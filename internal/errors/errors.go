package errors

// Пакет описания переменных ошибок
// Заменяет написание ошибок строками

import err "errors"

var ErrNotFound = err.New("not found")
var ErrAlreadyExist = err.New("already exist")
var ErrConfict = err.New("conflict")
var ErrLinked = err.New("linked")
var ErrInvalidID = err.New("invalid id")
