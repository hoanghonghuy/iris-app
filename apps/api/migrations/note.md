## install migrate CLI tool
- go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

## install postgresql client
- window: 
    - winget install PostgreSQL.PostgreSQL.16
    - or: https://www.postgresql.org/download/windows/ (just Command Line Tools)
    - or use docker: docker run --rm -it postgres:16-alpine psql [connection-string]
- linux: 
    - ubuntu/debian: sudo apt-get install postgresql-client
    - fedora/rhel: sudo dnf install postgresql
    - arch: sudo pacman -S postgresql-libs

## run migration

### dev
- migrate -path apps/api/migrations -database "postgres://postgres:iris@localhost:5433/iris_db?sslmode=disable" up
- migrate -path apps/api/migrations -database "postgres://postgres:iris@localhost:5433/iris_db?sslmode=disable" version

### production
- set env:
    - powershell: $env:DATABASE_URL="postgres://[user]:[password]@localhost:[port]/[database]?sslmode=disable"
    - linux: export DATABASE_URL="postgres://[user]:[password]@localhost:[port]/[database]?sslmode=disable"
- migrate -path apps/api/migrations -database "$env:DATABASE_URL" up
- migrate -path apps/api/migrations -database "$env:DATABASE_URL" down
- migrate -path apps/api/migrations -database "$env:DATABASE_URL" version

## seed database

### dev (local docker)
- docker exec -i iris-postgres psql -U postgres -d iris_db < scripts/db/seed_demo.sql
- docker exec -i iris-postgres psql -U postgres -d iris_db < scripts/db/seed_master.sql

### production
- cleanup data (keep schema):
    - powershell: psql $env:DATABASE_URL -f scripts/db/cleanup.sql
    - linux: psql $DATABASE_URL -f scripts/db/cleanup.sql
    - docker: docker run --rm -v ${PWD}/scripts/db:/scripts postgres:16-alpine psql $env:DATABASE_URL -f /scripts/cleanup.sql
- seed all data (1400+ rows):
    - powershell: psql $env:DATABASE_URL -f scripts/db/seed_master.sql
    - linux: psql $DATABASE_URL -f scripts/db/seed_master.sql
    - docker: docker run --rm -v ${PWD}/scripts/db:/scripts postgres:16-alpine psql $env:DATABASE_URL -f /scripts/seed_master.sql

## reset database (drop all + recreate)
- powershell: psql $env:DATABASE_URL -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
- linux: psql $DATABASE_URL -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
- docker: docker run --rm postgres:16-alpine psql $env:DATABASE_URL -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
- then run migrations: migrate -path apps/api/migrations -database "$env:DATABASE_URL" up
- then seed: psql $env:DATABASE_URL -f scripts/db/seed_master.sql
