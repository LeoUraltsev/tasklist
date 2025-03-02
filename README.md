# Task List Api

Запуск:   
1. Подготовить файл конфига .yaml в папке `./config`   
    Пример конфига:
    ```yaml
    app:
      env: "local"
    
    http:
      address: "localhost:8080"
    
    jwt:
      secret: "your_secret_key"
      exp: 1h
    
    storage:
      sqlite:
        path: "./storage/tasklist.db"
    ```

2. Запустить миграции

    ```bash
    make goose_up
    ```

3. Запустить приложение
    ```bash
    make run
    ```