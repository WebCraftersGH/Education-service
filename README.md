# Education-service
For LearnFlow.

## Что есть в сервисе

### Entities

- Problem (Задача из пулла задач)
  - ID
  - Name
  - Slug
  - Difficulty
  - Tag
  - AuthorID
  - VerifiedAt
  - CreatedAt
  - UpdatedAt

- ProblemContent (Контент задачи)
  - ID
  - ProblemID
  - DescriptionMD
  - InputFormatMD
  - OutputFormatMD
  - ConstraintsMD
  - NotesMD
  - CreatedAt
  - UpdatedAt

## Конфигурация

Сервис читает переменные из `.env` файла.

Пример:

```env
APP_ENV=local
HTTP_PORT=8080
LOG_LEVEL=info
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=education_service
DB_SSLMODE=disable
DB_TIMEZONE=UTC
```

Также в репозитории есть `.env.example`.

## Что делает при старте

1. Загружает конфиг из `.env`.
2. Инициализирует логгер.
3. Подключается к Postgres через GORM.
4. Выполняет `AutoMigrate` для таблиц `problems` и `problem_contents`.
5. Запускает HTTP сервер.

## Логи

Логи пишутся:
- в stdout
- в файл `logs/log.log`

Добавлены:
- стартовые логи приложения
- логирование HTTP запросов
- логи CRUD операций для `Problem` и `ProblemContent`
