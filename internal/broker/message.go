package broker

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Kuznecowilya/crypto-rates/internal/domain"
)

// RateMessage — формат курса в топике rates.raw.
type RateMessage struct {
	Name     string    `json:"name"`
	Currency string    `json:"currency"`
	Price    float64   `json:"price"`
	Source   string    `json:"source"`
	Time     time.Time `json:"time"`
	Stale    bool      `json:"stale"`
}

// Key возвращает ключ Kafka вида «BTC:USD».
func Key(name, currency string) string {
	return name + ":" + currency
}

// EncodeRate превращает курс в JSON для Kafka.
func EncodeRate(k domain.Kurs) ([]byte, error) {
	return json.Marshal(RateMessage{
		Name:     k.Name,
		Currency: k.Currency,
		Price:    k.Price,
		Source:   k.Source,
		Time:     k.TimeK,
		Stale:    k.Stale,
	})
}

// DecodeRate разбирает JSON и создаёт проверенный курс.
func DecodeRate(data []byte) (domain.Kurs, error) {
	var message RateMessage
	if err := json.Unmarshal(data, &message); err != nil {
		return domain.Kurs{}, fmt.Errorf("разбор сообщения курса: %w", err)
	}
	return domain.NewKurs(
		message.Name,
		message.Currency,
		message.Price,
		message.Source,
		message.Time,
		message.Stale,
	)
}
