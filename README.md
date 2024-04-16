# Сервис баннеров

## API

:rocket: Деплой API доступен по пути: https://milchenko.online/api/v2

Для поддерживания функционала с двумя ролями (токенами) было добавлено несколько методов в API. Их можно найти в [файле](https://drive.google.com/file/d/1QyiyX-Loq33jli2f-iuoQg4bFOUtoYR-/view?usp=sharing) или по [ссылке](https://www.postman.com/martian-eclipse-581372/workspace/bannersapi/collection/24084701-db52c204-fc41-4dde-a619-b3720ef0ad53?action=share&creator=24084701). Файл можно скачать для импорта в приложении.

## Как запустить проект

### Мануально

Можно запустить мануально по комаде, указанной ниже, но для этого понадобится поднятые базы Redis и PostgreSQL с настройками, указанными в конфиге по пути `configs/app/local.yml`

```sh
go run cmd/main.go -ConfigPath={путь_до_конфигурационного_файла}
```

Сами примеры файлов конфигурации можно посмотреть в `configs/app/local.yml`

### С помощью docker-compose

Здесь достаточно прописать в корне проекта:

```sh
docker-compose up -d
```

## Заполнение базы данных

Для заполенения базы данных можно воспользоваться скриптом, который находится по пути: `migrations/postgres_filling.py`
Как запускать:

```sh
python migrations/postgres_filling {url_до_API} {admin_token}
```

## E2E тесты

Для запуска тестов нужно прописать

```sh
pytest e2e/
```

## Makefile

Многие готовые сценарии уже обренуты в Makefile. Поэтому можно ознакомиться с ними, прописав:

```sh
make help
```
