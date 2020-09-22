# Chat API
[![Build Status](https://travis-ci.org/Filet-de-S/avito_chat.svg?branch=master)](https://travis-ci.org/Filet-de-S/avito_chat)

*Cервис – HTTP API для работы с чатами и сообщениями пользователя*
    
*Проект – [тестовое задание](https://github.com/Filet-de-S/avito_chat/blob/master/task.md) на позицию стажёра бэкенд-разработчика в **start.avito.unit.Messenger***

## Запуск сервиса
```bash
make prod
```
___
    
* [API Docs](https://app.swaggerhub.com/apis-docs/Filet-de-S/ChatAPI/1.0.0)

* Написаны интеграционные тесты, результат генерируется в `test/postman/result.html`

    * Для запуска в контейнере: `make test-run`
    * На локальной OS: `make test-local-run`

        Тестовый скрипт использует `postgresql-client, curl, newman, newman-reporter-html`, позаботьтесь об их установке

* Профилируем: перед запуском сервиса – `export PPROF=ON`
    
    URI: `/admin/pprof`
    
    Header: `Authorization: <secrets/.pprof>`

    Нагружаем и снимаем показатели с помощью готового скрипта `ApacheBench`: `make ab [args="[ab args]"]`
    
    * Если ругается, регулируйте время снятия показателей `make pprof tlim=[seconds]` и/или таймаут сервиса в `deployments/service.env`
    
    или с помощью `wrk`: `make wrk [tlim=[seconds] args="[wrk args]"` 

___
#### *NB: для удобства быстрой инициализации проекта оставил* 
* директорию `secrets/`: 

     `.pgpass` – PGPASSFILE
        
     `.pgpassf` – POSTGRES_PASSWORD_FILE
     
     `.pprof` – пароль для входа в /admin/pprof
        
     `.pwmng` – `user:pw` для входа в сервис-менеджер паролей
        
     `.uuids` – UUIDs v5 для генерации "уникальных" id. Необходимые поля: 
        `USER:uuid`, `CHAT:uuid`, `MSG:uuid`
        
     `migrations/002_create-role_client.up.sql` – создать роль клиента БД 
