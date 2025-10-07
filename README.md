# Image Processor

Асинхронный сервис обработки изображений с веб-интерфейсом. Позволяет загружать изображения, ставить задачи на обработку в очередь (Apache Kafka), и получать результаты в фоновом режиме.

## 🏗️ Архитектура

```
image-processor/
├── cmd/imageprocessor/          # Точка входа приложения
│   └── main.go
├── internal/                    # Внутренняя логика приложения
│   ├── api/                    # HTTP API слой
│   │   ├── handler/           # HTTP обработчики
│   │   │   ├── create.go      # POST /upload
│   │   │   ├── get.go         # GET /image/{id}, GET /image/info/{id}
│   │   │   ├── delete.go      # DELETE /image/{id}
│   │   │   ├── handler.go     # Структура Handler
│   │   │   └── interface.go   # Интерфейсы сервисов
│   │   ├── response/          # HTTP ответы
│   │   │   └── response.go
│   │   ├── router/            # Маршрутизация
│   │   │   └── router.go
│   │   └── server/           # HTTP сервер
│   │       └── server.go
│   ├── config/               # Конфигурация
│   │   ├── config.go
│   │   └── types.go
│   ├── kafka/               # Kafka клиент
│   │   └── kafka.go
│   ├── model/              # Модели данных
│   │   └── model.go
│   ├── repository/         # Слой данных
│   │   ├── minio/         # MinIO (файловое хранилище)
│   │   │   └── minio.go
│   │   └── postgres/      # PostgreSQL (метаданные)
│   │       ├── create.go
│   │       ├── delete.go
│   │       ├── get.go
│   │       ├── update.go
│   │       └── repo.go
│   └── service/           # Бизнес-логика
│       ├── create.go      # Создание изображений
│       ├── delete.go      # Удаление изображений
│       ├── get.go         # Получение изображений
│       ├── images.go      # Обработка изображений
│       ├── interface.go   # Интерфейсы
│       ├── service.go    # Основной сервис
│       ├── utils.go       # Утилиты обработки
│       └── workers.go     # Фоновые воркеры
├── migrations/            # SQL миграции
│   └── 20251006163803_message_table.sql
├── web/                  # Веб-интерфейс
│   ├── index.html        # Главная страница
│   ├── styles.css        # Стили
│   ├── app.js           # JavaScript логика
│   └── font.ttf         # Шрифт для водяных знаков
├── env/                 # Конфигурация
│   └── config.yaml
├── docker-compose.yml   # Docker окружение
├── Dockerfile          # Docker образ
├── go.mod             # Go модули
└── go.sum
```

## 🚀 Запуск через Docker

### 1. Клонирование и подготовка
```bash
git clone <repository-url>
cd image-processor
```

### 2. Запуск всех сервисов
```bash
docker-compose up --build -d
```

Это запустит:
- **PostgreSQL** (порт 5432) - база данных для метаданных
- **MinIO** (порт 9000) - файловое хранилище
- **Kafka + Zookeeper** - очередь сообщений
- **Image Processor** (порт 8080) - основное приложение

### 3. Проверка работы
```bash
# Проверить статус контейнеров
docker-compose ps

# Посмотреть логи
docker-compose logs -f image-processor
```

### 4. Доступ к сервисам
- **Веб-интерфейс**: http://localhost:8080
- **MinIO Console**: http://localhost:9001 (admin/admin123)
- **PostgreSQL**: localhost:5432 (postgres/postgres)

## 🌐 Веб-интерфейс

Откройте http://localhost:8080 для доступа к веб-интерфейсу.

### Возможности:
- **📤 Загрузка изображений** - drag & drop или выбор файла
- **⚙️ Выбор задач обработки**:
  - `resize` - изменение размера
  - `miniature generating` - создание миниатюр
  - `watermark` - добавление водяных знаков
- **📊 Мониторинг статуса** - реальное время обработки
- **🖼️ Предпросмотр результатов** - автоматическое отображение
- **💾 Скачивание файлов** - в папку загрузок браузера
- **🗑️ Удаление изображений** - с подтверждением

### Интерфейс:
- **Тёмная тема** с современным дизайном
- **Адаптивная верстка** для мобильных устройств
- **Уведомления** о статусе операций
- **Прогресс-бары** для отслеживания обработки

## 📡 API Endpoints

### POST /api/upload
Загрузка изображения на обработку.

**Content-Type**: `multipart/form-data`

**Параметры**:
- `image` (file) - файл изображения (PNG, JPEG, GIF)
- `metadata` (string) - JSON с параметрами задачи

**Пример metadata**:
```json
{
  "task": "resize",
  "content_type": "image/jpeg",
  "watermark_string": "© My Brand",
  "resize": {
    "width": 800,
    "height": 600
  }
}
```

**Ответ**:
```json
{
  "result": "550e8400-e29b-41d4-a716-446655440000"
}
```

### GET /api/image/{id}
Получение обработанного изображения.

**Параметры**:
- `id` (string) - UUID изображения

**Ответ**: Бинарные данные изображения или JSON с ошибкой

### GET /api/image/info/{id}
Получение информации о статусе обработки.

**Параметры**:
- `id` (string) - UUID изображения

**Ответ**:
```json
{
  "result": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "finished",
    "created_at": "2025-10-07T19:25:32Z"
  }
}
```

### DELETE /api/image/{id}
Удаление изображения.

**Параметры**:
- `id` (string) - UUID изображения

**Ответ**:
```json
{
  "result": "Image deleted successfully"
}
```

## 🔧 Поддерживаемые задачи

### 1. Resize (изменение размера)
```json
{
  "task": "resize",
  "resize": {
    "width": 800,
    "height": 600
  }
}
```

### 2. Thumbnail (миниатюра)
```json
{
  "task": "miniature generating"
}
```
Создаёт миниатюру 200x200 пикселей.

### 3. Watermark (водяной знак)
```json
{
  "task": "watermark",
  "watermark_string": "© My Brand"
}
```

## 🛠️ Технологии

- **Backend**: Go, Gin, PostgreSQL, MinIO, Apache Kafka
- **Frontend**: HTML5, CSS3, JavaScript (Vanilla)
- **Containerization**: Docker, Docker Compose
- **Image Processing**: Go image packages, OpenType fonts

## 📝 Конфигурация

Основные настройки в `env/config.yaml`:
- База данных PostgreSQL
- MinIO файловое хранилище
- Kafka брокер
- Директории для временных файлов

## 🔍 Мониторинг

- **Логи**: структурированные JSON логи
- **Метрики**: количество обработанных изображений
- **Статус**: HTTP health checks
- **Очередь**: мониторинг Kafka топиков

## 🚨 Обработка ошибок

- Валидация форматов изображений
- Проверка размеров файлов
- Обработка ошибок сети
- Retry механизмы для Kafka
- Graceful shutdown воркеров

## 📈 Производительность

- **Параллельная обработка**: 3 воркера по умолчанию
- **Асинхронная архитектура**: неблокирующие операции
- **Масштабируемость**: горизонтальное масштабирование воркеров
- **Кэширование**: временные файлы в памяти
