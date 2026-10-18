package domain

import (
	"errors"
	"math"
	"testing"
	"time"
)

func TestNewKurs_Normalize(t *testing.T) {
	msk := time.FixedZone("MSK", 3*60*60)
	inputTime := time.Date(2026, 1, 1, 12, 0, 0, 0, msk)

	k, err := NewKurs(
		"  btc ",
		" usd ",
		65000,
		" CoinGecko ",
		inputTime,
		true,
	)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	if k.Name != "BTC" {
		t.Errorf("Name = %q, ждали %q", k.Name, "BTC")
	}
	if k.Currency != "USD" {
		t.Errorf("Currency = %q, ждали %q", k.Currency, "USD")
	}
	if k.Source != "coingecko" {
		t.Errorf("Source = %q, ждали %q", k.Source, "coingecko")
	}
	if !k.Stale {
		t.Error("Stale = false, ждали true")
	}
	if k.TimeK.Location() != time.UTC {
		t.Errorf("время не в UTC: %v", k.TimeK.Location())
	}
	if k.TimeK.Hour() != 9 {
		t.Errorf("час = %d, ждали 9", k.TimeK.Hour())
	}
	if !k.TimeK.Equal(inputTime) {
		t.Errorf(
			"момент времени изменился: получили %v, был %v",
			k.TimeK,
			inputTime,
		)
	}
}

func TestNewKurs_KeepsPriceAndSource(t *testing.T) {
	k, err := NewKurs(
		"btc",
		"usd",
		65000.5,
		"coingecko",
		time.Now(),
		false,
	)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	if k.Price != 65000.5 {
		t.Errorf("Price = %v, ждали 65000.5", k.Price)
	}
	if k.Source != "coingecko" {
		t.Errorf("Source = %q, ждали %q", k.Source, "coingecko")
	}
	if k.Stale {
		t.Error("Stale = true, ждали false")
	}
}

func TestNewKurs_Invalid(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		coin     string
		currency string
		price    float64
		source   string
		at       time.Time
	}{
		{"отрицательная цена", "btc", "usd", -1, "coingecko", now},
		{"нулевая цена", "btc", "usd", 0, "coingecko", now},
		{"цена NaN", "btc", "usd", math.NaN(), "coingecko", now},
		{"бесконечная цена", "btc", "usd", math.Inf(1), "coingecko", now},
		{"пустое название", "", "usd", 100, "coingecko", now},
		{"название из пробелов", "   ", "usd", 100, "coingecko", now},
		{"пустая валюта", "btc", "", 100, "coingecko", now},
		{"пустой источник", "btc", "usd", 100, "", now},
		{"нулевое время", "btc", "usd", 100, "coingecko", time.Time{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			k, err := NewKurs(
				tc.coin,
				tc.currency,
				tc.price,
				tc.source,
				tc.at,
				false,
			)

			if err == nil {
				t.Fatalf("ждали ошибку, получили курс %+v", k)
			}
			if !errors.Is(err, ErrInvalidKurs) {
				t.Errorf(
					"ошибка не завёрнута в ErrInvalidKurs: %v",
					err,
				)
			}
			if k != (Kurs{}) {
				t.Errorf(
					"при ошибке ждали пустой Kurs, получили %+v",
					k,
				)
			}
		})
	}
}
