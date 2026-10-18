package domain

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// Kurs — один курс одной криптовалюты в конкретный момент времени.
type Kurs struct {
	Name     string
	Currency string
	Price    float64
	Source   string
	TimeK    time.Time
	Stale    bool
}

// NewKurs приводит входные данные к единому виду и не позволяет создать
// некорректный курс.
func NewKurs(
	name, currency string,
	price float64,
	source string,
	t time.Time,
	stale bool,
) (Kurs, error) {
	k := Kurs{
		Name:     strings.ToUpper(strings.TrimSpace(name)),
		Currency: strings.ToUpper(strings.TrimSpace(currency)),
		Price:    price,
		Source:   strings.ToLower(strings.TrimSpace(source)),
		TimeK:    t.UTC(),
		Stale:    stale,
	}
	if err := k.Validate(); err != nil {
		return Kurs{}, err
	}
	return k, nil
}

// Validate проверяет обязательные поля курса.
func (k Kurs) Validate() error {
	if strings.TrimSpace(k.Name) == "" {
		return fmt.Errorf("%w: нет названия монеты", ErrInvalidKurs)
	}
	if strings.TrimSpace(k.Currency) == "" {
		return fmt.Errorf("%w: нет данных о валюте", ErrInvalidKurs)
	}
	if k.Price <= 0 || math.IsNaN(k.Price) || math.IsInf(k.Price, 0) {
		return fmt.Errorf("%w: невозможная цена монеты %v", ErrInvalidKurs, k.Price)
	}
	if strings.TrimSpace(k.Source) == "" {
		return fmt.Errorf("%w: нет источника курса", ErrInvalidKurs)
	}
	if k.TimeK.IsZero() {
		return fmt.Errorf("%w: нулевое значение времени", ErrInvalidKurs)
	}
	return nil
}
