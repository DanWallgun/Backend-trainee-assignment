# Решение тестового задания на позицию стажера-бекендера в юнит Авто

## Запуск веб-сервера
Два контейнера: оффициальный docker образ MongoDB и контейнер, на котором работает сам api-веб-сервис на go. Первый образ подгружается, второй образ собирается локально по Dockerfile.  
Команда запуска:  
`docker-compose up --build`

## Overview
- Язык программирования: Go
- Для персистентности используется MongoDB
- Юнит-тесты только для http handlers. [Coverage report](https://htmlpreview.github.io/?https://github.com/DanWallgun/Backend-trainee-assignment/blob/master/api/pkg/handlers/test_handler_cover.html)
- Валидация URL
- Возможность задавать кастомные ссылки

## JSON API
- Сохранить короткое представление заданного URL
    - URL:  
    `/create`
    - Method:  
    `POST`
    - Параметры:  
    `long_url` - (обязательный) сокращаемый URL  
    `short_url` - (необязательный) пользовательское краткое представление; если опущено, будет выдано случайно
    ```
    {
        "long_url": "<URL>",
        "short_url": "<alphanumeric>"
    }
    ```
    - Response  
    `views` - количество переходов по сокращенной ссылке
    ```
    {
        "long_url": "<URL>",
        "short_url": "<alphanumeric>",
        "views": "<integer>"
    }
    ```
- Redirect на соответствующую страницу
    - URL:  
    `/{short_url}`
    - Method:  
    `GET`
- Информация о сокращении
    - URL:  
    `/info/{short_url}`
    - Method:  
    `GET`
    - Response аналогичен созданию сокращенной ссылки
## Пример запроса
```
curl  --header "Content-Type: application/json" \
      --request POST \
      --data '{"long_url":"https://github.com/avito-tech/auto-backend-trainee-assignment"}' \
      http://127.0.0.1:8080/create

{"long_url":"https://github.com/avito-tech/auto-backend-trainee-assignment","short_url":"c73k0b","views":0}
```  
```
curl --request GET http://127.0.0.1:8080/info/c73k0b

{"long_url":"https://github.com/avito-tech/auto-backend-trainee-assignment","short_url":"c73k0b","views":2}
```
