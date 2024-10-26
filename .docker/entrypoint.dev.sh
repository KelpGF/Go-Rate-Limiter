rm .env
cp .env.example .env

go mod tidy
go run main.go
