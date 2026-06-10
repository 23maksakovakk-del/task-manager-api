# Task Manager API – 7 реализаций на разных языках

Этот репозиторий содержит **7 полностью независимых реализаций** одного и того же REST API для управления задачами с JWT-аутентификацией и ролевой моделью (user/admin). Каждая версия демонстрирует один и тот же набор эндпоинтов:

- `POST /auth/register` – регистрация пользователя  
- `POST /auth/login` – получение JWT-токена  
- `GET /tasks?page=1&limit=10&status=pending` – список задач (с пагинацией и фильтром) с учётом роли: user видит только задачи, назначенные ему, admin – все задачи  
- `POST /tasks` – создание задачи  
- `PUT /tasks/:id` – обновление (только автор или admin)  
- `DELETE /tasks/:id` – удаление (только admin)

## Директории

| Папка | Язык / Фреймворк | База данных |
|-------|------------------|--------------|
| `01-nodejs-express` | Node.js + Express + Prisma | PostgreSQL |
| `02-python-fastapi` | Python + FastAPI + SQLAlchemy | PostgreSQL / SQLite |
| `03-go-gin` | Go + Gin + GORM | PostgreSQL |
| `04-java-spring` | Java + Spring Boot + JPA | PostgreSQL |
| `05-rust-actix` | Rust + Actix-web + SQLx | PostgreSQL |
| `06-csharp-dotnet` | C# + .NET 8 + EF Core | PostgreSQL |
| `07-php-laravel` | PHP + Laravel + Eloquent | PostgreSQL / MySQL |

## Как запустить

В каждой папке есть свой `README.md` с инструкцией. Общий порядок:

1. Установите PostgreSQL (или измените строку подключения)
2. Скопируйте `.env.example` в `.env` и укажите данные БД
3. Выполните миграции (в каждой реализации по-своему)
4. Запустите сервер (команда в каждой папке)
5. Используйте Postman или curl для тестирования

## Пример запроса

```bash
# Регистрация
curl -X POST http://localhost:3000/auth/register -H "Content-Type: application/json" -d '{"email":"test@ex.com","name":"Test","password":"123"}'

# Логин
curl -X POST http://localhost:3000/auth/login -H "Content-Type: application/json" -d '{"email":"test@ex.com","password":"123"}'

# Получение задач (токен из предыдущего шага)
curl -X GET "http://localhost:3000/tasks?page=1&limit=5" -H "Authorization: Bearer <token>"
