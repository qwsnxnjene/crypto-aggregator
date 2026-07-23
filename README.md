Агрегатор цен на популярные криптовалюты в реальном времени. Обновления по HTTP или WebSocket. Источники: Binance и CoinGecko.

## Как запустить
1. Создать БД и применить миграцию:
   psql -U postgres -d crypto_aggregator -f migrations/001_init.sql

2. Создать .env файл:
   DB_DSN=postgres://user:password@localhost:5432/table_name?sslmode=disable

3. Запустить бэкенд:
   go run cmd/main.go

4. Открыть frontend/index.html в браузере

## Архитектура
Пакет aggregator - сбор данных из разных источников для разных криптовалют
Пакет fetcher - логика получения цен от каждого конкретного источника
Пакет handler - обработка http-запросов и установления WebSocket-соединения
Пакет storage - работа с БД

## Технологии
Backend - Golang 1.25.11
Frontend - HTML + CSS
DB - Postgres
