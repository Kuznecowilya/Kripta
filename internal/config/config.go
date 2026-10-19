package config

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

type KafkaConfig struct {
	Brokers          []string
	RatesTopic       string
	PriceAlertsTopic string
	WriterGroup      string
	BotGroup         string
}

type FetcherConfig struct {
	CoinGeckoBaseURL string
	Symbols          []string
	Currencies       []string
	Interval         time.Duration
	Prefetch         time.Duration
	HTTPAddr         string
	Kafka            KafkaConfig
}

type WriterConfig struct {
	Postgres                 PostgresConfig
	Kafka                    KafkaConfig
	MinAlertThresholdPercent float64
	HTTPAddr                 string
}

type APIConfig struct {
	Postgres PostgresConfig
	HTTPAddr string
}

type BotConfig struct {
	Postgres      PostgresConfig
	Kafka         KafkaConfig
	APIBaseURL    string
	TelegramToken string
	HTTPAddr      string
	AlertCooldown time.Duration
}

func requiredEnv(name string) (string, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return "", fmt.Errorf("Пустое значение %s", name)
	}
	return value, nil
}

func requiredInt(name string) (int, error) {
	value, err := requiredEnv(name)
	if err != nil {
		return 0, err
	}
	number, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("Не вышло преобразовать строку %q из %s в число, ошибка: %w", value, name, err)
	}
	return number, nil
}

func requiredDuration(name string) (time.Duration, error) {
	value, err := requiredEnv(name)
	if err != nil {
		return 0, err
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("Не вышло преобразовать строку %q из %s в продолжительность времени, ошибка: %w", value, name, err)
	}
	return duration, nil
}

func requiredFloat(name string) (float64, error) {
	value, err := requiredEnv(name)
	if err != nil {
		return 0, err
	}
	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("Не вышло преобразовать строку %q из %s в дробное число, ошибка: %w", value, name, err)
	}
	if math.IsNaN(number) || math.IsInf(number, 0) {
		return 0, fmt.Errorf("Не вышло преобразовать строку %q из %s в дробное число, значение должно быть конечным числом", value, name)
	}
	return number, nil
}

func requiredList(name string) ([]string, error) {
	value, err := requiredEnv(name)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(value, ",")

	var result []string
	for i, part := range parts {
		item := strings.TrimSpace(part)
		if item == "" {
			return nil, fmt.Errorf("Значение №%d строки %q из %s пустое", i+1, value, name)
		}
		result = append(result, item)
	}
	return result, nil
}

func upperAll(items []string) []string {
	for i, item := range items {
		items[i] = strings.ToUpper(item)
	}
	return items
}

func LoadPostgresConfig() (PostgresConfig, error) {
	var problems []error
	host, hostErr := requiredEnv("POSTGRES_HOST")
	if hostErr != nil {
		problems = append(problems, hostErr)
	}
	port, portErr := requiredInt("POSTGRES_PORT")
	if portErr != nil {
		problems = append(problems, portErr)
	} else if port > 65535 || port < 1 {
		problems = append(problems, fmt.Errorf("Порт %d не соответствует стандарту", port))
	}
	user, userErr := requiredEnv("POSTGRES_USER")
	if userErr != nil {
		problems = append(problems, userErr)
	}
	password, passwordErr := requiredEnv("POSTGRES_PASSWORD")
	if passwordErr != nil {
		problems = append(problems, passwordErr)
	}
	database, databaseErr := requiredEnv("POSTGRES_DB")
	if databaseErr != nil {
		problems = append(problems, databaseErr)
	}
	sslMode := strings.TrimSpace(os.Getenv("POSTGRES_SSLMODE"))
	if sslMode == "" {
		sslMode = "disable"
	}
	if len(problems) > 0 {
		return PostgresConfig{}, errors.Join(problems...)
	}
	return PostgresConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
		SSLMode:  sslMode,
	}, nil
}

func LoadKafkaConfig() (KafkaConfig, error) {
	var problems []error
	brokers, brokersErr := requiredList("KAFKA_BROKERS")
	if brokersErr != nil {
		problems = append(problems, brokersErr)
	}
	ratesTopic, ratesTopicErr := requiredEnv("KAFKA_TOPIC_RATES_RAW")
	if ratesTopicErr != nil {
		problems = append(problems, ratesTopicErr)
	}
	priceAlertsTopic, priceAlertsTopicErr := requiredEnv("KAFKA_TOPIC_PRICE_ALERTS")
	if priceAlertsTopicErr != nil {
		problems = append(problems, priceAlertsTopicErr)
	}
	writerGroup, writerGroupErr := requiredEnv("KAFKA_GROUP_WRITER")
	if writerGroupErr != nil {
		problems = append(problems, writerGroupErr)
	}
	botGroup, botGroupErr := requiredEnv("KAFKA_GROUP_BOT")
	if botGroupErr != nil {
		problems = append(problems, botGroupErr)
	}

	if len(problems) > 0 {
		return KafkaConfig{}, errors.Join(problems...)
	}
	return KafkaConfig{
		Brokers:          brokers,
		RatesTopic:       ratesTopic,
		PriceAlertsTopic: priceAlertsTopic,
		WriterGroup:      writerGroup,
		BotGroup:         botGroup,
	}, nil
}

func LoadFetcherConfig() (FetcherConfig, error) {
	var problems []error

	coinGeckoBaseURL, coinGeckoBaseURLErr := requiredEnv("COINGECKO_BASE_URL")
	if coinGeckoBaseURLErr != nil {
		problems = append(problems, coinGeckoBaseURLErr)
	}

	symbols, symbolsErr := requiredList("CRYPTO_SYMBOLS")
	if symbolsErr != nil {
		problems = append(problems, symbolsErr)
	}
	symbols = upperAll(symbols)

	currencies, currenciesErr := requiredList("FIAT_CURRENCIES")
	if currenciesErr != nil {
		problems = append(problems, currenciesErr)
	}
	currencies = upperAll(currencies)

	interval, intervalErr := requiredDuration("FETCH_INTERVAL")
	if intervalErr != nil {
		problems = append(problems, intervalErr)
	} else if interval <= 0 {
		problems = append(problems, fmt.Errorf("FETCH_INTERVAL должен быть больше нуля, получено %s", interval))
	}

	prefetch, prefetchErr := requiredDuration("FETCH_PREFETCH")
	if prefetchErr != nil {
		problems = append(problems, prefetchErr)
	} else if prefetch <= 0 {
		problems = append(problems, fmt.Errorf("FETCH_PREFETCH должен быть больше нуля, получено %s", prefetch))
	}
	if intervalErr == nil && prefetchErr == nil && prefetch >= interval {
		problems = append(problems, fmt.Errorf("FETCH_PREFETCH (%s) должен быть меньше FETCH_INTERVAL (%s)", prefetch, interval))
	}

	httpAddr, httpAddrErr := requiredEnv("FETCHER_HTTP_ADDR")
	if httpAddrErr != nil {
		problems = append(problems, httpAddrErr)
	}

	kafka, kafkaErr := LoadKafkaConfig()
	if kafkaErr != nil {
		problems = append(problems, kafkaErr)
	}

	if len(problems) > 0 {
		return FetcherConfig{}, errors.Join(problems...)
	}
	return FetcherConfig{
		CoinGeckoBaseURL: coinGeckoBaseURL,
		Symbols:          symbols,
		Currencies:       currencies,
		Interval:         interval,
		Prefetch:         prefetch,
		HTTPAddr:         httpAddr,
		Kafka:            kafka,
	}, nil
}

func LoadAPIConfig() (APIConfig, error) {
	var problems []error

	postgres, postgresErr := LoadPostgresConfig()
	if postgresErr != nil {
		problems = append(problems, postgresErr)
	}

	httpAddr, httpAddrErr := requiredEnv("API_HTTP_ADDR")
	if httpAddrErr != nil {
		problems = append(problems, httpAddrErr)
	}

	if len(problems) > 0 {
		return APIConfig{}, errors.Join(problems...)
	}
	return APIConfig{
		Postgres: postgres,
		HTTPAddr: httpAddr,
	}, nil
}

func LoadWriterConfig() (WriterConfig, error) {
	var problems []error

	postgres, postgresErr := LoadPostgresConfig()
	if postgresErr != nil {
		problems = append(problems, postgresErr)
	}

	kafka, kafkaErr := LoadKafkaConfig()
	if kafkaErr != nil {
		problems = append(problems, kafkaErr)
	}

	threshold, thresholdErr := requiredFloat("ALERT_MIN_THRESHOLD_PERCENT")
	if thresholdErr != nil {
		problems = append(problems, thresholdErr)
	} else if threshold <= 0 {
		problems = append(problems, fmt.Errorf("ALERT_MIN_THRESHOLD_PERCENT должен быть больше нуля, получено %g", threshold))
	}

	httpAddr, httpAddrErr := requiredEnv("WRITER_HTTP_ADDR")
	if httpAddrErr != nil {
		problems = append(problems, httpAddrErr)
	}

	if len(problems) > 0 {
		return WriterConfig{}, errors.Join(problems...)
	}
	return WriterConfig{
		Postgres:                 postgres,
		Kafka:                    kafka,
		MinAlertThresholdPercent: threshold,
		HTTPAddr:                 httpAddr,
	}, nil
}
