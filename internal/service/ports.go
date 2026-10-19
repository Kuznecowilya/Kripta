package service

import (
	"context"

	"github.com/Kuznecowilya/crypto-rates/internal/domain"
)

type CryptoProvider interface {
	GetRates(ctx context.Context, symbols, currencies []string) ([]domain.Kurs, error)
}

type RatePublisher interface {
	Publish(ctx context.Context, kurs domain.Kurs) error
}
