db_login:
	psql ${DATABASE_URL}

db_create_migration:
	migrate create -ext sql -dir db/migrations -seq ${name}

db_migrate:
	migrate -database ${DATABASE_URL} -path db/migrations up
