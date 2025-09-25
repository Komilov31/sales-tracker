# Sales Tracker API

## Обзор

Sales Tracker - это веб-приложение для отслеживания финансовых операций, включая доходы (доход) и расходы (расход). Оно позволяет пользователям создавать, читать, обновлять и удалять записи с деталями, такими как сумма, дата и категория. Приложение предоставляет агрегированную аналитику (сумма, среднее, количество, медиана, 90-й процентиль) за диапазоны дат и поддерживает экспорт данных в формате CSV. Включен простой статический фронтенд (HTML/JS/CSS) для взаимодействия с пользователем, а API документирован с помощью Swagger.

Бэкенд построен на Go, обслуживает RESTful API с Gin. Данные хранятся в PostgreSQL. Проект использует чистую, многослойную архитектуру для разделения ответственности.

## Технологии

- **Язык**: Go
- **Веб-фреймворк**: Gin  для маршрутизации HTTP и middleware
- **База данных**: PostgreSQL
- **Логирование**: Zerolog
- **Валидация**: go-playground/validator
- **Документация API**: Swagger
- **Фронтенд**: HTML, CSS, JavaScript (без фреймворка)
- **Сборка/Развертывание**: Docker (Dockerfile, docker-compose.yml для Postgres + приложения)

## Архитектура

Проект следует многослойной архитектуре:

1. **Слой представления (Handlers)**: `internal/handler/` - Обработчики Gin для HTTP-запросов. Привязывают JSON, валидируют ввод, вызывают сервисы, возвращают ответы. Включает утилиты для общих операций и тесты.

2. **Слой бизнес-логики (Services)**: `internal/service/` - Основная логика для операций CRUD, расчетов аналитики (агрегации вроде суммы, среднего, медианы через SQL), генерации CSV/экспорта.

3. **Слой доступа к данным (Repositories)**: `internal/repository/` - Взаимодействия с базой данных с использованием PostgreSQL. Обрабатывает запросы для создания, чтения (все/фильтрованные), обновления, удаления, аналитики. Включает утилиты и специфические запросы аналитики.

4. **Модели/DTO**: `internal/model/` - Структуры вроде `Item` (ID, Type, Amount, Date, Category, CreatedAt, Aggregated) и `Aggregated` (Sum, Average, Count, Median, Percentile_90). `internal/dto/` для DTO запросов/ответов (например, CreateItem, UpdateItem).

5. **Конфигурация**: `internal/config/` - Загружает из `config/config.yaml` или переменных окружения. Типы для конфига Postgres и HTTP-сервера.

6. **Валидация**: `internal/validator/` - Пользовательские валидаторы для вводов.

7. **База данных**: PostgreSQL с миграциями в `migrations/` (например, create_message_table.sql – вероятно, для таблицы items). Соединение пулировано с макс 10 открытыми/5 простаивающими.

8. **Статические ассеты**: `static/` - index.html (основной UI с формами для добавления/просмотра/фильтрации/аналитики/экспорта), script.js (AJAX-вызовы к API), styles.css.

9. **Документы**: `docs/` - Автоматически генерируемые Swagger JSON/YAML из аннотаций.

Приложение работает на `localhost:8080` по умолчанию. Swagger UI на `/swagger/index.html`.

## Установка и настройка

### С использованием Docker (Рекомендуется)
1. Клонируйте репозиторий:
   ```
   git clone https://github.com/Komilov31/sales-tracker.git
   cd sales-tracker
   ```

2. Скопируйте `.env.example` в `.env` и обновите значения (например, учетные данные БД):
   ```
   cp .env.example .env
   ```

3. Запустите сервисы с Docker Compose (Postgres + приложение):
   ```
   docker-compose up -d
   ```


### Тестирование
- Юнит тесты: `go test ./...`
- Интеграция: Запустите приложение, используйте curl или Swagger для тестирования эндпоинтов.

## Эндпоинты API

Базовый URL: `http://localhost:8080`

Все эндпоинты используют JSON (кроме экспортов CSV). Аутентификация: Отсутствует (добавьте при необходимости). Ответы об ошибках: JSON с сообщениями (например, 400: {"error": "Invalid date"}).

### Страницы
- **GET /**  
  Получить основную HTML-страницу.  
  Curl:  
  ```
  curl -X GET http://localhost:8080/
  ```  
  Ответ: HTML-контент.

### Записи (CRUD)
Записи представляют доходы/расходы: Type ("доход" или "расход"), Amount (>0), Date (YYYY-MM-DD), Category (строка).

- **POST /items**  
  Создать новую запись.  
  Тело:  
  ```json
  {
    "type": "доход",
    "amount": 1000,
    "date": "2024-01-01",
    "category": "Зарплата"
  }
  ```  
  Curl:  
  ```
  curl -X POST http://localhost:8080/items \
    -H "Content-Type: application/json" \
    -d '{"type":"доход","amount":1000,"date":"2024-01-01","category":"Зарплата"}'
  ```  
  Ответ (200):  
  ```json
  {
    "id": 1,
    "type": "доход",
    "amount": 1000,
    "date": "2024-01-01",
    "category": "Зарплата",
    "created_at": "2024-01-01T00:00:00Z"
  }
  ```

- **GET /items**  
  Получить все записи (опционально sort_by: csv-список вроде "date,amount").  
  Curl:  
  ```
  curl -X GET "http://localhost:8080/items?sort_by=date&sort_by=amount"
  ```  
  Ответ (200): Массив записей (без агрегированных данных).

- **PUT /items/{id}**  
  Обновить запись по ID (частичные обновления).  
  Тело: например, `{"amount": 1500}`  
  Curl:  
  ```
  curl -X PUT http://localhost:8080/items/1 \
    -H "Content-Type: application/json" \
    -d '{"amount":1500}'
  ```  
  Ответ (200): `{"message": "Item updated successfully"}`

- **DELETE /items/{id}**  
  Удалить запись по ID.  
  Curl:  
  ```
  curl -X DELETE http://localhost:8080/items/1
  ```  
  Ответ (200): `{"message": "Item deleted successfully"}`

- **GET /items/csv**  
  Экспортировать записи как CSV (опционально sort_by).  
  Curl:  
  ```
  curl -X GET "http://localhost:8080/items/csv?sort_by=date" \
    --output items.csv
  ```  
  Ответ: Скачивание CSV-файла.

### Аналитика
Агрегированные статистики по категории/типу за диапазон дат (from/to: YYYY-MM-DD). Включает сумму, среднее, количество, медиану, 90-й процентиль.

- **GET /analytics**  
  Получить агрегированные записи.  
  Curl:  
  ```
  curl -X GET "http://localhost:8080/analytics?from=2024-01-01&to=2024-12-31"
  ```  
  Ответ (200): Массив записей с `aggregated_data`.

- **GET /analytics/csv**  
  Экспортировать аналитику как CSV.  
  Curl:  
  ```
  curl -X GET "http://localhost:8080/analytics/csv?from=2024-01-01&to=2024-12-31" \
    --output analytics.csv
  ```  
  Ответ: CSV-файл.

### Документация Swagger
- **GET /swagger/*any**  
  Доступ к Swagger UI.  
  Curl (или браузер):  
  ```
  curl -X GET http://localhost:8080/swagger/index.html
  ```  
  Откройте `http://localhost:8080/swagger/index.html` в браузере для интерактивной документации.

