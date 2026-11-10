# crypto-rates — план обучения

Один урок рассчитан примерно на 4–5 часов. Отметки:

- `[ ]` — не начат;
- `[~]` — начат;
- `[x]` — завершён и проверен.

Тесты пишет агент. Интеграционные проверки выполняются дома.

## Блок 1. Ядро и контракты

### [~] Урок 1. Alert, Subscription и Stats

**Цель:** описать доменные данные уведомления, подписки и статистики, не
привязывая их к Kafka, SQL, HTTP или Telegram.

**Файлы:**

- `internal/domain/alert.go`;
- `internal/domain/subscription.go`;
- `internal/domain/stats.go`;
- тесты пишет агент.

**Готово, когда:** конструкторы и `Validate` не позволяют создать неверные
значения, `MaxChange` выбирает изменение с наибольшим модулем, тесты проходят.

### [ ] Урок 2. Все interfaces в ports.go

**Цель:** описать договоры между бизнес-логикой и внешним миром.

**Файл:** `internal/service/ports.go`.

**Добавить:** `AlertPublisher`, `RateRepository`,
`SubscriptionRepository`, `Notifier`, `RatesAPI`, `Clock`, `SystemClock`.

## Блок 2. Бизнес-логика

### [ ] Урок 3. FetchRates: сбор и fallback

**Цель:** независимо собирать пары и использовать предварительный курс с
`Stale=true`, если основной курс пары недоступен.

**Файл:** `internal/service/fetch_rates.go`.

**Методы:** `Collect`, `Publish`, helpers пар.

### [ ] Урок 4. FetchRates: цикл времени

**Цель:** немедленный первый запуск, затем работа по границам интервала с одним
prefetch за 5 секунд.

**Методы:** `Cycle`, `Loop`, `waitUntil`.

### [ ] Урок 5. IngestRates

**Цель:** сохранять курс, рассчитывать изменения USD за 1/2/3 часа и
публиковать alert при достижении порога.

**Файл:** `internal/service/ingest_rates.go`.

### [ ] Урок 6. RatesQuery

**Цель:** получить последние курсы, историю и статистику; проверить параметры.

**Файл:** `internal/service/rates_query.go`.

### [ ] Урок 7. NotifyAlerts и BotCommands

**Цель:** фильтровать уведомления по порогу и cooldown; обрабатывать команды
`/start`, `/help`, `/rates`, `/subscribe`, `/unsubscribe`.

**Файлы:**

- `internal/service/notify_alerts.go`;
- `internal/service/bot_commands.go`.

## Блок 3. HTTP и общая обвязка

### [ ] Урок 8. CoinGecko client

**Цель:** реализовать `CryptoProvider` через `/simple/price`.

**Файл:** `internal/exchange/coingecko.go`.

### [ ] Урок 9. Logger, server, app и API client

**Цель:** общие логи, `/health`, graceful shutdown, параллельный запуск задач и
HTTP-клиент бота к своему API.

**Файлы:**

- `internal/logger/logger.go`;
- `internal/httpserver/server.go`;
- `internal/app/app.go`;
- `internal/apiclient/client.go`.

## Блок 4. Kafka

### [ ] Урок 10. Kafka messages и producers

**Цель:** закончить JSON-сообщения курсов и alerts, написать producers.

**Файлы:**

- `internal/broker/message.go` — часть курса уже готова;
- `internal/broker/producer.go`.

### [ ] Урок 11. Kafka consumer и topics

**Цель:** ручной commit после обработки, retry временных ошибок, `ErrSkip` для
битых сообщений, создание топиков при старте.

**Файлы:**

- `internal/broker/consumer.go`;
- `internal/broker/topics.go`.

## Блок 5. PostgreSQL

### [ ] Урок 12. Миграции и подключение

**Цель:** таблицы `rates`, `subscriptions`, `deliveries`, ожидание базы и
безопасные параллельные миграции.

**Файлы:**

- `internal/storage/postgres/migrations/001_init.sql`;
- `internal/storage/postgres/postgres.go`.

### [ ] Урок 13. RateRepository

**Цель:** реализовать `Save`, `PriceNear`, `Latest`, `LatestAll`, `History`.

**Файл:** `internal/storage/postgres/repo.go`.

### [ ] Урок 14. SubscriptionRepository

**Цель:** подписки, отключение, список активных подписок, cooldown и
дедупликация доставок.

**Файл:** `internal/storage/postgres/repo.go`.

## Блок 6. Первая рабочая цепочка

### [ ] Урок 15. cmd/fetcher и cmd/writer

**Цель:** связать config, CoinGecko, Kafka, PostgreSQL и service use cases.

**Файлы:**

- `cmd/fetcher/main.go`;
- `cmd/writer/main.go`.

### [ ] Урок 16. Первый домашний запуск

**Цель:** поднять Docker, выполнить тесты и проверить:

```text
CoinGecko -> fetcher -> Kafka -> writer -> PostgreSQL
```

## Блок 7. REST и Telegram

### [ ] Урок 17. REST handlers

**Цель:** маршруты курсов, истории и статистики, коды ошибок, logging и recover
middleware.

**Файл:** `internal/transport/httpapi/api.go`.

### [ ] Урок 18. cmd/api и проверка REST

**Цель:** собрать API и проверить endpoints через `curl`.

**Файл:** `cmd/api/main.go`.

### [ ] Урок 19. Telegram HTTP client

**Цель:** `sendMessage`, `getUpdates`, JSON и защита token от попадания в
сетевые ошибки.

**Файл:** `internal/transport/telegram/telegram.go`.

### [ ] Урок 20. Telegram polling и cmd/bot

**Цель:** long polling, команды, Kafka alerts и health server в одном процессе.

**Файл:** `cmd/bot/main.go`.

## Блок 8. Сквозная проверка

### [ ] Урок 21. Настоящий Telegram и end-to-end

**Цель:** проверить все четыре процесса и полный путь сообщения до телефона.

### [ ] Урок 22. Финальная проверка и Git

**Цель:** `gofmt`, `go vet`, `go test -race`, graceful shutdown, отсутствие
секретов, чистый Git и push.

## Отдельные поздние уроки

Они не входят в основной рабочий MVP:

- transactional outbox;
- Prometheus metrics;
- Grafana;
- Alloy и Loki;
- Dockerfile для приложений;
- OpenAPI;
- CI.
