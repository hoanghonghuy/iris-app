-- install migrate CLI tool
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

-- seed database
docker exec -i iris-postgres psql -U postgres -d iris_db < scripts/db/seed_demo.sql

-- run migration
migrate -path apps/api/migrations -database "postgres://postgres:iris@localhost:5433/iris_db?sslmode=disable" up
migrate -path apps/api/migrations -database "postgres://postgres:iris@localhost:5433/iris_db?sslmode=disable" version
