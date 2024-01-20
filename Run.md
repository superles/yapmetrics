# Запуск агента
go run -ldflags "-X main.buildVersion=v1.0.0" .\cmd\agent\main.go

# Запуск сервера
go run -ldflags "-X main.buildVersion=v1.0.0" .\cmd\server\main.go