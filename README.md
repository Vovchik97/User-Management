# Система управления пользователями

API-сервер для управления пользователями, группами и ролями.

## 📚 Технологии

- Go 1.24
- Gin Web Framework
- GORM (PostgreSQL)
- JWT для аутентификации
- Swagger для документации API
- Docker и Docker Compose для контейнеризации

## 🚀 Запуск с использованием Docker Compose

### 📋 Предварительные требования

- [Docker](https://www.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### ⚙️ Шаги для запуска

1. Клонируйте репозиторий:

```bash
git clone https://github.com/Vovchik97/User-Management.git
cd userManagement
```

2. Соберите приложение с помощью Docker Compose:
```bash
docker-compose build
```

3. Запустите приложение с помощью Docker Compose:

```bash
docker-compose up -d
```

Эта команда запустит два контейнера:
- API-сервер на порту 8080
- PostgreSQL на порту 5432

При первом запуске автоматически создаётся учётная запись администратора:
- Email: admin@example.com
- Пароль: admin123

4.Проверьте работу API:

Откройте в браузере: http://localhost:8080/docs или http://localhost:8080/swagger/index/html

### 🛑 Остановка приложения

Если хочешь просто остановить, но чтобы контейнер остался в Docker Desktop:

```bash
docker-compose stop
```

Полная остановка и удаление контейнеров, но без уаления данных БД:

```bash
docker-compose down
```

Для удаления всех данных, включая тома базы данных:

```bash
docker-compose down -v
```

## ⚙️ Переменные окружения

Все необходимые переменные окружения уже настроены в файле `docker-compose.yml`. При необходимости вы можете изменить их значения:

- `DB_HOST` - хост базы данных (postgres)
- `DB_USER` - пользователь базы данных (postgres)
- `DB_PASSWORD` - пароль базы данных (password)
- `DB_NAME` - имя базы данных (user_management)
- `DB_PORT` - порт базы данных (5432)
- `JWT_SECRET` - секретный ключ для JWT токенов

Также пример настроек находится в файле .env.example.

## 📖 Документация API

После запуска приложения документация API доступна по адресу:
- http://localhost:8080/docs
- http://localhost:8080/swagger/index/html

## 📡 Основные маршруты API

| Метод | Путь | Описание |
|:------|:-----|:---------|
| `POST` | `/register` | Регистрация нового пользователя |
| `POST` | `/login` | Аутентификация и получение JWT токена |
| `GET` | `/users` | Получение списка всех пользователей (для админов и модераторов) |
| `GET` | `/users/:id` | Получение информации о пользователе по ID |
| `PUT` | `/users/:id` | Обновление данных пользователя |
| `DELETE` | `/users/:id` | Удаление пользователя (только для админа) |
| `POST` | `/groups` | Создание новой группы |
| `GET` | `/groups` | Получение списка групп |
| `GET` | `/groups/:id` | Получение информации о группе по ID |
| `PUT` | `/groups/:id` | Обновление данных группы |
| `DELETE` | `/groups/:id` | Удаление группы |
| `POST` | `/groups/:id/users/:userId` | Добавление пользователя в группу |
| `DELETE` | `/groups/:id/users/:userId` | Удаление пользователя из группы |
| `POST` | `/roles` | Создание новой роли |
| `GET` | `/roles` | Получение списка всех ролей |
| `POST` | `/users/:id/assign-role` | Назначение роли пользователю |
| `GET` | `/logs` | Просмотр логов действий пользователей (для админа) |
| `GET` | `/docs` | Swagger-документация API |