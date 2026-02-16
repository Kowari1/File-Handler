# Поднять контейнер
docker-compose up -d

# Сделать миграцию
goose -dir migrations postgres host=127. 0. 0. 1 port=5435 user=postgres password=123 dbname=devices sslmode=disable up

# Запустить приложение
go run cmd/app/main. go

# REST API
http: //localhost: 8080/web/ - страница сайта

GET http: //localhost: 8080/devices? page=1&limit=10 - девайсы с пагинацией
GET http: //localhost: 8080/devices/{guid}? page=1&limit=5 - конкретный дeвайc по guid
