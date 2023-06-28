# Cистема отслеживания городских мероприятий

## Описание

Этот проект - API системы для отслеживания городских мероприятий.
Этот API позовляет пользователям искать различные мероприятия, а так же регистрироваться на них.

### Используемые технологии

- Golang
    - Fiber - http фреймворк.
- PostgreSQL - База данных
- Swagger - документация API
- JWT

## Запуск

### Подготовка

Для начала склонируйте проект

```bash
git clone https://github.com/burenotti/rtu-it-lab-recruit
```

Затем нужно создать приватный RSA ключ. Он нужен для выдачи и валидации JWT токенов.

```bash
openssl genrsa -des3 -out private.pem 3072
openssl rsa -in private.pem -out private_unencrypted.pem -outform PEM
```

Теперь перейдем к настройке окружения.

Для начала создайте файл `.env.db`. В нем будут храниться переменные среды,
связанные с контейнером PostgreSQL.

```dotenv
# .env.db
POSTGRES_PASSWORD=secretpassword
POSGRES_DB=ems
```

После этого настроим само приложение.

Нам нужно указать две вещи: как подключиться к базе и как подключиться к smtp серверу.
Почтовый (SMTP) сервер необходим для активации пользователей. Вход в систему так же происходит через почту.
Вы можете создать тестовый сервер на [Mailtrap](https://mailtrap.io). Или можно просто
написать [мне](https://t.me/burenotti).
Я выдам вам доступ к своему серверу :)

Создайте в корне проекта файл `.env`.

```dotenv
# .env

DB_DSN=postgres://postgres:secretpassword@db:5432/ems

# Здесь настройки вашего сервера.
SMTP_HOST=
SMTP_PORT=
SMTP_USER=
SMTP_PASSWORD=
```

Теперь мы готовы к запуску.

Поднимите боевой сервер с помощью этой команды.

```bash
make prod # Пока что эта команда не существует...
```

Ура, мы в эфире!

Чтобы проверить, что все работает зайдите в [Swagger документацию](http://localhost:8080/docs)
