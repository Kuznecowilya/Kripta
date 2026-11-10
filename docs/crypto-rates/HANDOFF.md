# crypto-rates — HANDOFF для локального и Cloud Agent

Этот файл — источник актуального контекста проекта. Новый агент должен прочитать
его вместе с `LESSONS.md` до того, как давать следующее задание.

## 1. Цель проекта

Учебный сервис:

1. получает курсы `BTC/USD`, `BTC/RUB`, `ETH/USD`, `ETH/RUB` из CoinGecko;
2. передаёт их через Kafka;
3. сохраняет историю в PostgreSQL;
4. отдаёт курсы, историю и статистику через REST API;
5. отправляет уведомления через Telegram при заметном изменении цены.

Технологии:

- Go 1.25;
- PostgreSQL 16;
- Kafka в KRaft-режиме;
- REST API;
- Telegram Bot API;
- Docker Compose.

Поздние дополнительные этапы: Prometheus, Grafana, Alloy, Loki,
transactional outbox, CI и Dockerfile приложений.

## 2. Архитектура

Четыре запускаемых процесса:

```text
CoinGecko -> fetcher -> Kafka rates.raw
Kafka rates.raw -> writer -> PostgreSQL
                         -> Kafka price.alerts
client/bot -> REST API -> PostgreSQL
Telegram <-> bot -> REST API
Telegram <- bot <- Kafka price.alerts
```

### Fetcher

- делает один предварительный запрос за `FETCH_PREFETCH=5s`;
- на целевой границе времени делает основной запрос;
- публикует каждую пару независимо;
- если основная пара не получена, использует предварительную с `Stale=true`;
- не подключается к PostgreSQL.

### Writer

- читает `rates.raw`;
- повторно валидирует курс;
- сохраняет его с `ON CONFLICT DO NOTHING`;
- для свежих USD-курсов считает изменение за 1, 2 и 3 часа;
- при достижении минимального порога публикует `PriceAlert`;
- подтверждает Kafka offset только после успешной обработки.

### API

- читает PostgreSQL;
- отдаёт последние курсы, историю и статистику;
- не читает Kafka и не знает Telegram.

### Bot

- получает курсы только через REST API;
- хранит свои подписки и доставки в PostgreSQL;
- читает `price.alerts`;
- учитывает индивидуальный порог и cooldown 3 часа;
- подавляет повтор события по `(event_id, chat_id)`.

## 3. Принятые технические решения

- Модуль: `github.com/Kuznecowilya/crypto-rates`.
- Kafka topics: `rates.raw`, `price.alerts`.
- Ключ сообщения: `BTC:USD`.
- Consumer groups: `rates-writer`, `telegram-price-alerts`.
- Alert считается только для USD.
- Окна alert: 1, 2 и 3 часа.
- Допуск поиска старой цены: 10 минут в обе стороны.
- Для раннего этапа alert публикуется сразу после сохранения курса.
  Transactional outbox — отдельный поздний этап.
- Время хранится и передаётся в UTC.
- Реальный `.env` не попадает в Git.

## 4. Профиль ученика

Шкала первоначальной самооценки: 0 — не знает, 3 — пишет самостоятельно.

Уверенно:

- переменные, `if`, `for`, функции;
- пакеты и импорты;
- базовые структуры и срезы;
- `strconv.Atoi`;
- основы SQL;
- общее понимание Kafka.

Нужны подробные короткие памятки:

- указатели;
- интерфейсы-порты;
- `context.Context`;
- горутины и каналы;
- `defer`;
- HTTP-клиент и HTTP-сервер;
- JSON;
- `time.Duration`;
- `pgx`, пул соединений и миграции;
- Kafka в Go, consumer groups и commit offset;
- Docker;
- Telegram Bot API;
- чистая архитектура.

Правила общения находятся в `.cursor/rules/crypto-rates.mdc`.

## 5. Работа на двух компьютерах

### Рабочий компьютер

- Windows;
- проект: `C:\Users\support9.TK\projects\crypto-rates`;
- Git установлен;
- код пишется в Cursor;
- Docker и полноценный запуск откладываются до дома.

### Домашний компьютер

- WSL, пользователь `demon`;
- проект: `~/projects/crypto-rates`;
- здесь запускаются Go, тесты, Docker, Kafka и PostgreSQL.

GitHub — единый источник файлов между компьютерами.

Рабочий цикл:

```text
перед работой: git pull
после работы:  git add . -> git commit -> git push
```

Секретный `.env` создаётся отдельно на каждом компьютере и не переносится через
GitHub.

## 6. Текущая ветка

Текущая учебная ветка:

```text
cursor/work-code-tests-ec3d
```

Она основана на:

```text
cursor/home-config-transfer-ec3d
```

Готовая эталонная реализация для сверки находится в ветке:

```text
cursor/crypto-rates-reference-ec3d
```

и в её папке `reference/`.

Эталон не нужно механически копировать без понимания. Пользователь собирает
свой проект по урокам.

## 7. Что уже есть в текущем проекте

### Domain

- `internal/domain/rate.go` — `Kurs`, `NewKurs`, `Validate`;
- `internal/domain/errors.go` — доменные sentinel-ошибки;
- `internal/domain/price_change.go` — `PercentChange`, `ReachedThreshold`,
  `ChangeDirection`;
- табличные тесты курса и процентов.

### Config

- `internal/config/config.go`;
- структуры `PostgresConfig`, `KafkaConfig`, `FetcherConfig`, `WriterConfig`,
  `APIConfig`, `BotConfig`;
- helpers `requiredEnv`, `requiredInt`, `requiredDuration`, `requiredFloat`,
  `requiredList`, `upperAll`;
- все шесть `Load*Config`;
- один prefetch (`FETCH_PREFETCH`);
- табличные тесты helpers и loaders;
- тест, что пароль и Telegram-токен не попадают в ошибку.

### Service

В `internal/service/ports.go` пока только:

- `CryptoProvider`;
- `RatePublisher`.

Остальные порты добавляются после доменных типов первого урока.

### Broker

В `internal/broker/message.go` готовы:

- `RateMessage`;
- `Key`;
- `EncodeRate`;
- `DecodeRate`.

Есть тесты round-trip, невалидных сообщений и Kafka key.

Часть для `PriceAlert` будет добавлена после `domain.PriceAlert`.

### Infrastructure

- `.env.example`;
- `.gitignore`;
- `docker-compose.yml` с PostgreSQL, Kafka и Kafka UI;
- `cmd/fetcher/main.go` пока содержит только компилируемый каркас.

## 8. Проверенное состояние

На ветке `cursor/work-code-tests-ec3d` проходят:

```text
go vet ./...
go test -race -count=1 ./...
```

При этом полный рабочий проект ещё не собран. Проверки относятся только к
написанной части.

## 9. Текущий урок

Текущий этап: **Урок 1 — доменные модели alert, subscription и stats**.

На рабочем компьютере создан локальный файл:

```text
internal/domain/alert.go
```

На момент последнего обновления он ещё не был закоммичен и находился только на
Windows-компьютере пользователя. Агент не должен перезаписывать его. Сначала
нужно прочитать текущий файл и проверить попытку пользователя.

Задание, которое было выдано:

1. добавить `package domain`;
2. добавить нужные imports;
3. описать `Change`;
4. описать `PriceAlert`.

После проверки нужно по частям написать:

- `AlertEventID`;
- `NewPriceAlert`;
- `PriceAlert.Validate`;
- `PriceAlert.MaxChange`.

Затем:

- `internal/domain/subscription.go`;
- `internal/domain/stats.go`;
- тесты к ним пишет агент.

## 10. Следующий шаг для нового локального агента

1. Прочитать локальный `internal/domain/alert.go`.
2. Не менять файл без запроса пользователя.
3. Проверить написанные структуры и imports.
4. Дать только следующее небольшое задание урока.
5. После завершения Урока 1 написать табличные тесты, запустить проверки,
   обновить этот HANDOFF и отметить урок в `LESSONS.md`.

## 11. Безопасность

- Не коммитить `.env`.
- Не коммитить реальный пароль PostgreSQL.
- Не коммитить Telegram bot token.
- Не просить вставлять секреты в чат.
- Не логировать URL Telegram, содержащий token.
