### Для запуска приложения в терминале 
Последовательность действий для тестирования приложения в терминале:

1. Запуск контейнера PostgresQL
    ```sh
    docker compose up -d db
    ```
2. Запуск в терминале команды
    ```sh
    cd backend && go run . 10, 11, 14, 15
    ```

### Для запуска приложения в браузере 

1. необходимо запустить контейнеры frontend (vue.js), backend и db (postgres):  
    ```bash
    docker compose up -d
    ```
2. открыть страницу http://localhost:5173/  
   
Инициализация базы данных происходит автоматически при старте контейнера. Для инициализации базы данных используется:  
 [схема](https://github.com/dave502/Go.Vue/blob/58996118444538bb72d798745ff4406ed1568128/backend/db/init/1.schema.sql)  
[скрипт заполнения данными](https://github.com/dave502/Go.Vue/blob/58996118444538bb72d798745ff4406ed1568128/backend/db/init/2.inserts.sql)
 