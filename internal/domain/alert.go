package domain

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type Change struct {
	Hours int
	Percent float64
}
type PriceAlert struct {
	EventID string
	Name string
	Currency string
	Price float64
	Changes []Change
	TimeK time.Time
}

func AlertEventID(k Kurs) string {
	return fmt.Sprintf("%s:%s:%d", (strings.ToUpper(strings.TrimSpace(k.Name))), (strings.ToUpper(strings.TrimSpace(k.Currency))), k.TimeK.UTC().Unix())
}

func NewPriceAlert(k Kurs, changes []Change) (PriceAlert, error){
	result := PriceAlert{
		EventID: AlertEventID(k),
		Name: (strings.ToUpper(strings.TrimSpace(k.Name))),
		Currency: (strings.ToUpper(strings.TrimSpace(k.Currency))),
		Price: k.Price,
		Changes: append([]Change(nil), changes...),
		TimeK: k.TimeK.UTC(),
		}	
		if err := result.Validate(); err != nil {
		return PriceAlert{}, err
	}
return result, nil
}
func (a PriceAlert) Validate() error{
	if strings.TrimSpace(a.EventID) == "" {
		return fmt.Errorf("%w: нет ID", ErrInvalidAlert)
	}
	if strings.TrimSpace(a.Name) == "" {
		return fmt.Errorf("%w: нет названия монеты", ErrInvalidAlert)
	}
	if strings.TrimSpace(a.Currency) == "" {
		return fmt.Errorf("%w: нет данных о валюте", ErrInvalidAlert)
	}
	if a.Price <= 0 || math.IsNaN(a.Price) || math.IsInf(a.Price, 0) {
		return fmt.Errorf("%w: невозможная цена монеты %v", ErrInvalidAlert, a.Price)
	}
	if len(a.Changes) == 0 {
		return fmt.Errorf("%w: нет данных об изменении цены", ErrInvalidAlert)
	}
	for _, c := range a.Changes {
		if c.Hours <= 0 {
			return fmt.Errorf("%w: некорректное количество часов %d", ErrInvalidAlert, c.Hours)
		}
		if math.IsNaN(c.Percent) || math.IsInf(c.Percent, 0){
			return fmt.Errorf("%w: некорректный процент изменения %v", ErrInvalidAlert, c.Percent)
		}
	}
	if a.TimeK.IsZero() {
		return fmt.Errorf("%w: нулевое значение времени", ErrInvalidAlert)
	}
	return nil
}
func (a PriceAlert) MaxChange() Change{
	
}