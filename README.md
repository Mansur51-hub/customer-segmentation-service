# Customer segmentation service

Микросервис для динамического сегментирования пользователей

Используемые технологии:
- PostgreSQL (в качестве хранилища данных)
- Docker (для запуска сервиса)
- Swagger (для документации API)
- Mux (веб фреймворк)
- viper (для конфигурации)
- pgx (драйвер для работы с PostgreSQL)

# Usage
Конфигурацию можно настроить в файле .env

Запустить сервис можно с помощью команды `docker compose up`

Миграция бд отсутствует. Конечный файл для создания всех необходимых таблиц находится в папке migrations

Документацию после завпуска сервиса можно посмотреть по адресу `http://localhost:8080/swagger/index.html`
с портом 8080 по умолчанию

## Examples

Некоторые примеры запросов
- [Добавление сегмента](#create-segment)
- [Удаление сегмента](#delete-segment)
- [Добавление пользователя в сегмент](#add-user)
- [Получение активных сегентов пользователя](#get-segments)
- [Получение списка операций в формате csv](#get-ops)

### Добавление сегмента <a name="create-segment"></a>

Добавление сегмента с заданным процентом пользователей:

```curl
curl -X 'POST' \
  'http://localhost:8080/segments' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "percent": 70,
  "slug": "AVITO_VOICE_MESSAGES"
}'
```

Добавление сегмента без автоматического заполнения:

```curl
curl -X 'POST' \
  'http://localhost:8080/segments' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "slug": "AVITO_DISCOUNT_50"
}'
```

Пример ответа:
```json
{
  "id": 3,
  "slug": "AVITO_DISCOUNT_50",
  "percent": 0
}
```

### Удаление сегмента <a name="delete-segment"></a>

Удаление сегмента:

```curl
curl -X 'DELETE' \
  'http://localhost:8080/segments' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "slug": "AVITO_VOICE_MESSAGES"
}'
```
Пример ответа:
```json
"success"
```

### Добавление пользователя в сегменты <a name="add-user"></a>

NOTE: для назначения ttl необходимо добавить поле "ttl" в формате time.Duration string 

Пример: 

```curl
curl -X 'POST' \
  'http://localhost:8080/users' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "segments_to_add": [
    {
      "slug": "AVITO_DISCOUNT_50"
    },
    {
     "slug": "AVITO_VOICE_MESSAGES",
      "ttl": "24h"  
    },
    {
     "slug": "AVITO_PERFORMANCE_VAS",
      "ttl": "48h"  
    }
  ],
  "segments_to_delete": [
    "AVITO_DISCOUNT_50"
  ],
  "user_id": 1000
}'

```
Пример ответа:
```json
[
  {
    "id": 1,
    "user_id": 1000,
    "segment_slug": "AVITO_DISCOUNT_50",
    "created_at": "2023-08-31T21:08:07.8915698+03:00",
    "expired_at": "0001-01-01T00:00:00Z"
  },
  {
    "id": 2,
    "user_id": 1000,
    "segment_slug": "AVITO_VOICE_MESSAGES",
    "created_at": "2023-08-31T21:08:07.8981185+03:00",
    "expired_at": "2023-09-01T21:08:07.898118Z"
  },
  {
    "id": 3,
    "user_id": 1000,
    "segment_slug": "AVITO_PERFORMANCE_VAS",
    "created_at": "2023-08-31T21:08:07.9009127+03:00",
    "expired_at": "2023-09-02T21:08:07.900912Z"
  }
]
```

### Получение активных сегентов пользователя <a name="get-segments"></a>

NOTE: запросы GET c body могут не выполняться в сваггере.

Запрос с curl из консоли выполняется всегда.

Пример:

```curl
curl -X 'GET' \
  'http://localhost:8080/users' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "limit": 10,
  "offset": 0,
  "user_id": 1000
}'
```
Пример ответа, с указанием id резервирования:
```json
["AVITO_PERFORMANCE_VAS","AVITO_VOICE_MESSAGES"]
```

### Получение списка операций в формате csv <a name="get-ops"></a>

Пример:

```curl
curl -X 'GET' \
  'http://localhost:8080/operations' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  --data-raw '{
  "limit": 10,
  "month": 8,
  "offset": 0,
  "year": 2023
}'

```

Пример ответа:
```json
1000,AVITO_DISCOUNT_50,adding,Thu Aug 31 21:08:07 2023
1000,AVITO_VOICE_MESSAGES,adding,Thu Aug 31 21:08:07 2023
1000,AVITO_PERFORMANCE_VAS,adding,Thu Aug 31 21:08:07 2023
1000,AVITO_DISCOUNT_50,removing,Thu Aug 31 21:08:07 2023
```
