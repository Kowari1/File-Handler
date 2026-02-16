docker-compose up -d
goose -dir migrations postgres "host=127.0.0.1 port=5435 user=postgres password=123 dbname=devices sslmode=disable" up
go run cmd/app/main.go
