package domain

import "errors"

// ошибка: курс не прошёл проверку
var ErrInvalidKurs = errors.New("invalid kurs")

// ошибка: неправильный интервал подписки
var ErrInvalidSubscription = errors.New("invalid subscription")

// ошибка: неправильная настройка экстренного уведомления
var ErrInvalidAlert = errors.New("invalid alert")

// ошибка: неизвестная или недопустимая валюта
var ErrInvalidCurrency = errors.New("invalid currency")

// ошибка: неизвестный или недопустимый часовой пояс
var ErrInvalidTimezone = errors.New("invalid timezone")

// ошибка: в базе не найдена нужная запись
var ErrNotFound = errors.New("not found")

// ошибка: монета отсутствует в списке разрешённых
var ErrUnknownSymbol = errors.New("unknown symbol")

// ошибка: старая цена равна нулю, процент вычислить нельзя
var ErrDivisionByZero = errors.New("division by zero")
