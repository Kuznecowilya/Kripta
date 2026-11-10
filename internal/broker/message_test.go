package broker

import (
	"testing"
	"time"

	"github.com/Kuznecowilya/crypto-rates/internal/domain"
)

func TestRateMessageRoundTrip(t *testing.T) {
	want, err := domain.NewKurs(
		"btc",
		"usd",
		65000.5,
		"coingecko",
		time.Date(2026, 10, 2, 12, 0, 0, 0, time.UTC),
		true,
	)
	if err != nil {
		t.Fatal(err)
	}

	data, err := EncodeRate(want)
	if err != nil {
		t.Fatal(err)
	}
	got, err := DecodeRate(data)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("получили %+v, ожидали %+v", got, want)
	}
}

func TestDecodeRateInvalid(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{"не JSON", `{`},
		{"нулевая цена", `{"name":"BTC","currency":"USD","price":0,"source":"x","time":"2026-10-02T12:00:00Z"}`},
		{"нет времени", `{"name":"BTC","currency":"USD","price":1,"source":"x"}`},
		{"нет названия", `{"currency":"USD","price":1,"source":"x","time":"2026-10-02T12:00:00Z"}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := DecodeRate([]byte(test.data)); err == nil {
				t.Fatal("ожидали ошибку")
			}
		})
	}
}

func TestKey(t *testing.T) {
	if got := Key("BTC", "USD"); got != "BTC:USD" {
		t.Errorf("Key() = %q, ожидали BTC:USD", got)
	}
}
