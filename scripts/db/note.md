docker exec -i iris-postgres psql -U postgres -d iris_db < scripts/db/cleanup.sql
docker exec -i iris-postgres psql -U postgres -d iris_db < scripts/db/seed_demo.sql