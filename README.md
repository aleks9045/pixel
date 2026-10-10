docker compose up --build

Деплоинг будет тяжелым, придется везде изменять url: в гитхабе, в nginx кофигурациях, в конфигах и .env

frontend и backend - отдельные проекты со своей документацией

на уровне с compose.yml должен быть файл .env:

```
POSTGRES_DB=postgres
POSTGRES_USER=postgres
```
