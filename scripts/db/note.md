## linux
docker exec -i iris-postgres psql -U postgres -d iris_db < scripts/db/cleanup.sql
docker exec -i iris-postgres psql -U postgres -d iris_db < scripts/db/seed_demo.sql

## deployment
- test connection
docker run --rm postgres:16 psql "$env:DATABASE_URL" -c "select now();"
- seed database
docker run --rm -v "$($PWD.Path):/work" postgres:16 psql "$env:DATABASE_URL" -f /work/scripts/db/seed_demo.sql

