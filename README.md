# API- Сервис с базовым Working Pool

## Дисклеймер перед началом:
- Я так понимаю, что тут надо было включить свою фантазию и постараться показать примитив работы с горутинами и каналами,
- Поэтому я решил создать АПИшку,  хотел еще Доккер и БД прикрепить (какую-нибудь  postgreSQL),  но
- не придумал что туда написать,  useless немного. Зато прикрутил `Swagger` для четкой документации.
- Ну а дальше буду в `README.md` описывать что я сделал.

# Стек:
- Golang
- Swagger
## Условие: 
Реализовать простейший примитив Worker Pool, где данные динамически добавляются и удаляются, воркеры принимают и работают
с данными (строками) затем высвобождаются спустя время. Задание на работу с `горутинами` и `каналами`.

## Бизнес логика:
Есть некоторый сервис который умеет принимать некоторые `ручки`.

Все ручки начинаются с префикса `api/`

### `api/ping`
Нужен для того, чтобы проверить что сервис функционирует и запускается. Должен по этой ручке ответить PONG и 
статус код 200

```yaml
GET /api/ping

Response:

     200 OK
    
Body:  
  
     PONG
```

### `api/add_worker`
Нужен для того, чтобы создать новый воркер динамически и новым последовательным ID (прошлый воркер 3 - этот 4)

```yaml
POST /api/add_worker

Response:

     200 OK

Body:
     
```

### `api/add_job`
Нужен для того, чтобы занять воркер работой (принятием строки). Необходимо передавать параметр `job` для корректной работы
пример строки: `api/add_job?job=ReadmeDa` если воркеров нет, то выдает ошибку об этом и не назначает работу

```yaml
POST /api/add_job?job=ReadmeDa

Response:

     200 OK

Body:
     
```

### `api/remove_worker`
Нужен для того чтобы удалить воркер, если воркеров нет, то выдает ошибку и не дает удалить
# Структура проекта
Весь проект разбит на файлы.
### `cmd/myapp/`
По этому пути расположен файл `main.go`, который содержит в себе все [ручки] и все самое главное для правильной работы проекта.

### `/internal/config/`
По этому пути расположен файл `config.go`, в котором находится функция, запускающая всю конфигурацию, (.env).

### `docs/`
По этому пути расположены файлы, которые отвечают за `Swagger`, для лучшего представления микросервиса.
### `/internal/api/`
По этому пути расположены файлы, которые отвечают за основную работу приложения (ручки).
- `handlers.go` отвечает за описание всех действий, связанных с воркер пулом.

### `/internal/worker/`
По этому пути расположены файлы, которые которые реализуют логику и структуру worker pool

### `validators/`
По этому пути расположен файл `Validate.go`, который отвечает за проверку при создании тендера (правильный ввод данных, правильная обработка их).

### `cmd/myapp/test/`
По этому пути расположены файлы для тестирования корректной функциональности приложения.
