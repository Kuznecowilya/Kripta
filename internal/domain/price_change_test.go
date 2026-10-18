package domain

import (
	"errors"
	"math"
	"testing"
)

func TestPercentChange(t *testing.T) {
	const epsilon = 0.000000001

	tests := []struct {
		name    string
		old     float64
		new     float64
		want    float64
		wantErr bool
	}{
		{"рост на 5 процентов", 100, 105, 5, false},
		{"падение на 5 процентов", 100, 95, -5, false},
		{"цена не изменилась", 100, 100, 0, false},
		{"падение BTC", 63000, 60000, -4.761904761904762, false},
		{"старая цена равна нулю", 0, 100, 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := PercentChange(tc.old, tc.new)

			if tc.wantErr {
				if !errors.Is(err, ErrDivisionByZero) {
					t.Fatalf("ошибка = %v, ждали ErrDivisionByZero", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}

			if math.Abs(got-tc.want) > epsilon {
				t.Errorf("результат = %v, ждали %v", got, tc.want)
			}
		})
	}
}

func TestReachedThreshold(t *testing.T) {
	tests := []struct {
		name      string
		change    float64
		threshold float64
		want      bool
	}{
		{"рост выше порога", 5, 3, true},
		{"падение выше порога", -5, 3, true},
		{"изменение равно порогу", 3, 3, true},
		{"изменение меньше порога", 2.99, 3, false},
		{"нулевой порог", 5, 0, false},
		{"отрицательный порог", 5, -1, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ReachedThreshold(tc.change, tc.threshold)

			if got != tc.want {
				t.Errorf("результат = %v, ждали %v", got, tc.want)
			}
		})
	}
}

func TestChangeDirection(t *testing.T) {
	tests := []struct {
		name   string
		change float64
		want   string
	}{
		{"цена выросла", 5, "up"},
		{"цена упала", -5, "down"},
		{"цена не изменилась", 0, "flat"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ChangeDirection(tc.change)

			if got != tc.want {
				t.Errorf("результат = %q, ждали %q", got, tc.want)
			}
		})
	}
}
