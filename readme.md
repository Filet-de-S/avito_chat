# Chat service
[![Build Status](https://travis-ci.org/Filet-de-S/avito_chat.svg?branch=master)](https://travis-ci.org/Filet-de-S/avito_chat)

**О проекте:**
    чат-сервис, предоставляющий HTTP API для работы с чатами и сообщениями пользователя.
    
*Проект – [тестовое задание](https://github.com/Filet-de-S/avito_chat/blob/master/task.md) на позицию стажёра-бэкенд разработчика в **start.avito.unit.Messenger***

## Запуск сервиса
```bash
make prod
```
___
    
* [API Docs](https://app.swaggerhub.com/apis-docs/Filet-de-S/ChatAPI/1.0.0)

* Написаны интеграционные тесты, результат генерируется в `test/postman/result.html`

    * Для запуска в контейнере: `make test-run`
    * На локальной OS: `make test-local-run`. 
    Тестовый скрипт использует `postgresql-client, curl, newman, newman-reporter-html`. Позаботьтесь об их наличии  


___
#### *NB: для удобства быстрой инициализации проекта оставил* 
* директорию `secrets/`: 

     `.pgpass` – PGPASSFILE
        
     `.pgpassf` – POSTGRES_PASSWORD_FILE
        
     `.pwmng` – `user:pw` для входа в сервис-менеджер паролей
        
     `.uuids` – UUIDs v5 для генерации "уникальных" id. Необходимые поля: 
        `USER:uuid`, `CHAT:uuid`, `MSG:uuid`
        
     `migrations/002_create-role_client.up.sql` – создать роль клиента БД 
