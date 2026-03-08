-- install migrate CLI tool
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

-- run migration
migrate -path apps/api/migrations -database "postgres://postgres:iris@localhost:5433/iris_db?sslmode=disable" up
migrate -path apps/api/migrations -database "postgres://postgres:iris@localhost:5433/iris_db?sslmode=disable" version
