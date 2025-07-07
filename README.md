# Wallet API

Простой REST API для пополнения и снятия средств с виртуальных кошельков.  
Реализован на Go с использованием PostgreSQL и Docker.

---

## 🔧 Функционал

- Создание кошелька (автоматически при первом обращении)
- Пополнение баланса (`DEPOSIT`)
- Списание средств (`WITHDRAW`)
- Получение текущего баланса по ID кошелька
- Полностью атомарные операции
- Работает в высоконагруженной среде (1000 RPS)

---

## 🛠 Технологии

- **Go 1.22** — основной язык
- **PostgreSQL** — хранение данных
- **Docker / Docker Compose** — контейнеризация
- **pgx/v4** — драйвер PostgreSQL
- **chi** — маршрутизатор HTTP
- **Clean Architecture** — разделение слоёв

---

## 📁 Структура проекта
```
/wallet-api
├── cmd
│ └── main.go # Точка входа
├── config
│ ├── config.go # Загрузка конфига
├── internal
│ ├── api # HTTP обработчики
│ ├── domain # Доменные модели и типы
│ ├── service # Бизнес-логика
│ └── repository # Работа с БД через pgx
├── migrations # SQL миграции (не используется — создаётся автоматически)
├── Dockerfile
├── docker-compose.yml
├── config.env # Конфиг окружения
└── README.md
```

---

## 🚀 Установка и запуск

### 1. Клонируй репозиторий

```bash
git clone https://github.com/Kuzahka/wallet-api.git 
cd wallet-api
```

## Собрать и запустить все сервисы
```
docker-compose --env-file config.env up --build
```

## Проверить статус контейнеров
```
docker-compose ps
```

## Посмотреть логи
```
docker-compose logs -f wallet-api
docker-compose logs -f postgres
```

## Остановить всё
```
docker-compose down
```

## Пополнение баланса (DEPOSIT)
```
POST /api/v1/wallet
Content-Type: application/json
```
```
{
  "walletId": "a1b2c3d4-e5f6-7890-g1h2-i3j4k5l6m7n8",
  "operationType": "DEPOSIT",
  "amount": 1000
}
```

## Списание средств (WITHDRAW)
```
POST /api/v1/wallet
Content-Type: application/json
```
```
{
  "walletId": "a1b2c3d4-e5f6-7890-g1h2-i3j4k5l6m7n8",
  "operationType": "WITHDRAW",
  "amount": 500
}
```
## Получить баланс кошелька
```
GET /api/v1/wallets/a1b2c3d4-e5f6-7890-g1h2-i3j4k5l6m7n8
```

## Пример ответа
```
{
  "id": "a1b2c3d4-e5f6-7890-g1h2-i3j4k5l6m7n8",
  "balance": 1000
}
```
