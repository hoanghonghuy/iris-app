## install migrate CLI tool
- go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

## seed database
- docker exec -i iris-postgres psql -U postgres -d iris_db < scripts/db/seed_demo.sql

## run migration

### dev
- migrate -path apps/api/migrations -database "postgres://postgres:iris@localhost:5433/iris_db?sslmode=disable" up
- migrate -path apps/api/migrations -database "postgres://postgres:iris@localhost:5433/iris_db?sslmode=disable" version

### production
- set env:
    powershell: $env:DATABASE_URL="postgres://[user]:[password]@localhost:[port]/[database]?sslmode=disable"
    linux: export DATABASE_URL="postgres://[user]:[password]@localhost:[port]/[database]?sslmode=disable"
- migrate -path apps/api/migrations -database "$env:DATABASE_URL" up
- migrate -path apps/api/migrations -database "$env:DATABASE_URL" down
- migrate -path apps/api/migrations -database "$env:DATABASE_URL" version