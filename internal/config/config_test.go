package config

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func validEnv() map[string]string {
	return map[string]string{
		"POSTGRES_HOST": "localhost", "POSTGRES_PORT": "5432",
		"POSTGRES_USER": "crypto", "POSTGRES_PASSWORD": "secret",
		"POSTGRES_DB": "crypto_rates", "POSTGRES_SSLMODE": "",
		"KAFKA_BROKERS":            "localhost:9092, localhost:9093",
		"KAFKA_TOPIC_RATES_RAW":    "rates.raw",
		"KAFKA_TOPIC_PRICE_ALERTS": "price.alerts",
		"KAFKA_GROUP_WRITER":       "rates-writer",
		"KAFKA_GROUP_BOT":          "telegram-price-alerts",
		"COINGECKO_BASE_URL":       "https://api.coingecko.com/api/v3",
		"CRYPTO_SYMBOLS":           "btc, Eth", "FIAT_CURRENCIES": "usd,rub",
		"FETCH_INTERVAL": "5m", "FETCH_PREFETCH": "5s",
		"ALERT_MIN_THRESHOLD_PERCENT": "3", "ALERT_COOLDOWN": "3h",
		"API_HTTP_ADDR": ":8080", "FETCHER_HTTP_ADDR": ":8082",
		"WRITER_HTTP_ADDR": ":8083", "BOT_HTTP_ADDR": ":8084",
		"API_BASE_URL": "http://localhost:8080", "TELEGRAM_BOT_TOKEN": "token",
	}
}

func setEnv(t *testing.T, overrides map[string]string) {
	t.Helper()
	env := validEnv()
	for name, value := range overrides {
		env[name] = value
	}
	for name, value := range env {
		t.Setenv(name, value)
	}
}

func TestConfigHelpers(t *testing.T) {
	t.Run("requiredEnv", func(t *testing.T) {
		tests := []struct {
			name, value, want string
			wantErr           bool
		}{
			{"значение", " abc ", "abc", false},
			{"пусто", "", "", true},
			{"пробелы", "   ", "", true},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				t.Setenv("TEST_VAR", test.value)
				got, err := requiredEnv("TEST_VAR")
				if got != test.want || (err != nil) != test.wantErr {
					t.Fatalf("получили %q, err=%v", got, err)
				}
			})
		}
	})

	t.Run("requiredInt", func(t *testing.T) {
		tests := []struct {
			value string
			want  int
			err   bool
		}{
			{"5432", 5432, false}, {"-1", -1, false},
			{"abc", 0, true}, {"1.5", 0, true}, {"", 0, true},
		}
		for _, test := range tests {
			t.Run(test.value, func(t *testing.T) {
				t.Setenv("TEST_VAR", test.value)
				got, err := requiredInt("TEST_VAR")
				if got != test.want || (err != nil) != test.err {
					t.Fatalf("получили %d, err=%v", got, err)
				}
			})
		}
	})

	t.Run("requiredDuration", func(t *testing.T) {
		tests := []struct {
			value string
			want  time.Duration
			err   bool
		}{
			{"10s", 10 * time.Second, false},
			{"1h30m", 90 * time.Minute, false},
			{"5", 0, true}, {"abc", 0, true}, {"", 0, true},
		}
		for _, test := range tests {
			t.Run(test.value, func(t *testing.T) {
				t.Setenv("TEST_VAR", test.value)
				got, err := requiredDuration("TEST_VAR")
				if got != test.want || (err != nil) != test.err {
					t.Fatalf("получили %v, err=%v", got, err)
				}
			})
		}
	})

	t.Run("requiredFloat", func(t *testing.T) {
		tests := []struct {
			value string
			want  float64
			err   bool
		}{
			{"2.5", 2.5, false}, {"0", 0, false},
			{"NaN", 0, true}, {"+Inf", 0, true}, {"abc", 0, true},
		}
		for _, test := range tests {
			t.Run(test.value, func(t *testing.T) {
				t.Setenv("TEST_VAR", test.value)
				got, err := requiredFloat("TEST_VAR")
				if got != test.want || (err != nil) != test.err {
					t.Fatalf("получили %v, err=%v", got, err)
				}
			})
		}
	})

	t.Run("requiredList", func(t *testing.T) {
		tests := []struct {
			value string
			want  []string
			err   bool
		}{
			{" BTC , ETH ", []string{"BTC", "ETH"}, false},
			{"BTC,,ETH", nil, true}, {"BTC,", nil, true}, {"", nil, true},
		}
		for _, test := range tests {
			t.Run(test.value, func(t *testing.T) {
				t.Setenv("TEST_VAR", test.value)
				got, err := requiredList("TEST_VAR")
				if !reflect.DeepEqual(got, test.want) || (err != nil) != test.err {
					t.Fatalf("получили %v, err=%v", got, err)
				}
			})
		}
	})

	t.Run("upperAll", func(t *testing.T) {
		got := upperAll([]string{"btc", "Eth", "USD"})
		want := []string{"BTC", "ETH", "USD"}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("получили %v, ожидали %v", got, want)
		}
	})
}

func TestLoadersValid(t *testing.T) {
	setEnv(t, nil)

	tests := []struct {
		name string
		load func() error
	}{
		{"Postgres", func() error { _, err := LoadPostgresConfig(); return err }},
		{"Kafka", func() error { _, err := LoadKafkaConfig(); return err }},
		{"Fetcher", func() error {
			got, err := LoadFetcherConfig()
			if err == nil && (!reflect.DeepEqual(got.Symbols, []string{"BTC", "ETH"}) ||
				!reflect.DeepEqual(got.Currencies, []string{"USD", "RUB"})) {
				t.Errorf("регистр не приведён: %+v", got)
			}
			return err
		}},
		{"Writer", func() error { _, err := LoadWriterConfig(); return err }},
		{"API", func() error { _, err := LoadAPIConfig(); return err }},
		{"Bot", func() error { _, err := LoadBotConfig(); return err }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.load(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestLoadersInvalid(t *testing.T) {
	tests := []struct {
		name      string
		overrides map[string]string
		load      func() error
		contains  []string
	}{
		{"Postgres собирает ошибки",
			map[string]string{"POSTGRES_HOST": "", "POSTGRES_PORT": "70000", "POSTGRES_DB": ""},
			func() error { _, err := LoadPostgresConfig(); return err },
			[]string{"POSTGRES_HOST", "Порт 70000", "POSTGRES_DB"}},
		{"Kafka собирает ошибки",
			map[string]string{"KAFKA_BROKERS": "", "KAFKA_GROUP_BOT": ""},
			func() error { _, err := LoadKafkaConfig(); return err },
			[]string{"KAFKA_BROKERS", "KAFKA_GROUP_BOT"}},
		{"Fetcher проверяет интервалы",
			map[string]string{"FETCH_INTERVAL": "5m", "FETCH_PREFETCH": "5m"},
			func() error { _, err := LoadFetcherConfig(); return err },
			[]string{"должен быть меньше FETCH_INTERVAL"}},
		{"Writer проверяет порог",
			map[string]string{"ALERT_MIN_THRESHOLD_PERCENT": "NaN"},
			func() error { _, err := LoadWriterConfig(); return err },
			[]string{"конечным числом"}},
		{"API собирает ошибки",
			map[string]string{"POSTGRES_HOST": "", "API_HTTP_ADDR": ""},
			func() error { _, err := LoadAPIConfig(); return err },
			[]string{"POSTGRES_HOST", "API_HTTP_ADDR"}},
		{"Bot собирает ошибки",
			map[string]string{"TELEGRAM_BOT_TOKEN": "", "ALERT_COOLDOWN": "0s"},
			func() error { _, err := LoadBotConfig(); return err },
			[]string{"TELEGRAM_BOT_TOKEN", "ALERT_COOLDOWN"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setEnv(t, test.overrides)
			err := test.load()
			if err == nil {
				t.Fatal("ожидали ошибку")
			}
			for _, text := range test.contains {
				if !strings.Contains(err.Error(), text) {
					t.Errorf("в ошибке нет %q:\n%v", text, err)
				}
			}
		})
	}
}

func TestErrorsDoNotLeakSecrets(t *testing.T) {
	setEnv(t, map[string]string{
		"POSTGRES_PORT":      "abc",
		"POSTGRES_PASSWORD":  "super-secret-password",
		"TELEGRAM_BOT_TOKEN": "super-secret-token",
	})
	_, err := LoadBotConfig()
	if err == nil {
		t.Fatal("ожидали ошибку")
	}
	for _, secret := range []string{"super-secret-password", "super-secret-token"} {
		if strings.Contains(err.Error(), secret) {
			t.Errorf("секрет %q попал в ошибку", secret)
		}
	}
}
