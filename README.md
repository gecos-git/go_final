# ToDo List
Этот репозиторий содержит REST API сервис Task Manager, который позволяет пользователям управлять заданиями и отслеживать их выполнение. Минимальный функциональный вариант в рамках дипломного проекта факультета Go-разработки, Яндекс Практикум.
# Эндпоинты

| URL              | HTTP      | Функциональность                |
| ---------------- | --------- | ------------------------------- |
| '/'              | GET       | домашняя страница               |
| '/api/task'      | POST      | создание задачи                 |
| '/api/task'      | GET       | получение задачи                |
| '/api/task'      | PUT       | редактирование задачи           |
| '/api/task'      | DELETE    | удаление задачи                 |
| '/api/tasks'     | GET       | получение списка задач       |
| '/api/nextdate'  | GET       | получение следующей даты задачи |
| '/api/task/done' | POST      | фиксация выполнения задачи      |

# Таблица БД

| Столбец    | Тип          | Описание                      |
| ---------- | ------------ | ----------------------------- |
| id         | INTEGER      | уникальный идентификатор      |
| date       | VARCHAR(8)   | дата задачи                   |
| title      | TEXT         | заголовок задачи              |
| comment    | TEXT         | комментарий к задаче          |
| repeat     | VARCHAR(128) | правила повторений для задачи |

# Запуск Сервиса

1. Клонируем репозитарий:
```shell
git clone git@github.com:gecos-git/go_final.git

cd go_final_project
```
2. Подтягиваем зависимости:
```shell
go mod download
```
3. Cобираем сервис:
```shell
make build
```
4. Запускаем сервис:
```shell
make run
```
5. Запускаем тесты:
```shell
make test


Адрес запущенного приложения:
http://localhost:7540/

Выход из приложения: ```Ctrl+C```