# Song API

API для работы с песнями и их текстами. Позволяет создавать, получать, обновлять и удалять песни, а также управлять их текстами.
Сразу предупреждаючто проект сделан с помощью нейросетей!
## Технологии

- Go 1.24+
- Gin (веб-фреймворк)
- GORM (ORM для работы с базой данных)
- PostgreSQL
- Swagger (документация API)

## Структура проекта

```
golang_test/
├── config/
│   └── database.go    # Конфигурация подключения к базе данных
├── handlers/
│   └── song_handler.go # Обработчики HTTP запросов
├── models/
│   ├── song.go        # Модели данных
│   └── song_detail.go # Модель детальной информации о песне
├── services/
│   └── song_service.go # Бизнес-логика
├── logger/
│   └── logger.go      # Логирование
├── docs/              # Swagger документация
├── main.go           # Точка входа
└── go.mod            # Зависимости проекта
```

## API Endpoints

### Песни

#### Получение списка песен
```
GET /api/songs
Query параметры:
- group (опционально) - фильтр по группе
- song (опционально) - фильтр по названию песни
- page (по умолчанию: 1) - номер страницы
- page_size (по умолчанию: 10) - размер страницы
```

#### Создание песни
```
POST /api/songs
Body:
{
    "group": "Название группы",
    "song": "Название песни"
}
```

#### Обновление песни
```
PUT /api/songs/:id
Body:
{
    "group": "Новое название группы",
    "song": "Новое название песни"
}
```

#### Удаление песни
```
DELETE /api/songs/:id
```

#### Получение информации о песне
```
GET /api/songs/info
Query параметры:
- group (обязательно) - название группы
- song (обязательно) - название песни
```

### Тексты песен

#### Получение текста песни
```
GET /api/songs/:id/lyrics
Query параметры:
- song_id (обязательно) - ID песни
- page (по умолчанию: 1) - номер страницы
- page_size (по умолчанию: 5) - размер страницы
```

## Примеры использования

### Создание песни
```bash
curl -X POST http://localhost:8080/api/songs \
  -H "Content-Type: application/json" \
  -d '{
    "group": "Muse",
    "song": "Supermassive Black Hole"
  }'
```

### Получение списка песен
```bash
curl "http://localhost:8080/api/songs?group=Muse&page=1&page_size=10"
```

### Получение информации о песне
```bash
curl "http://localhost:8080/api/songs/info?group=Muse&song=Supermassive%20Black%20Hole"
```

## Swagger документация

Swagger UI доступен по адресу:
```
http://localhost:8080/api/swagger/index.html
```

## Установка и запуск

1. Клонируйте репозиторий:
```bash
git clone https://github.com/EXTRAS12/golang_test.git
cd golang_test
```

2. Установите зависимости:
```bash
go mod download
```

3. Создайте файл .env с настройками базы данных:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=your_database
PORT=8080
```

4. Запустите приложение:
```bash
go run main.go
```

## База данных

### Таблица songs
- id (uint, primary key)
- group (string)
- song (string)

### Таблица lyrics
- id (uint, primary key)
- song_id (uint, foreign key)
- text (string)
- order (int) 
