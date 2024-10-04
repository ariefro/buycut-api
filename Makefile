ifneq (,$(wildcard .env.local))
	include .env.local
	export
endif

.PHONY: dbstart
dbstart: ## start the database server
	docker compose -f docker-compose.yml --env-file ./.env.local up --build

.PHONY: dbstop
dbstop: ## stop the database server
	docker compose -f docker-compose.yml --env-file ./.env.local down -v

.PHONY: createdb
createdb: ## create the database
	docker exec -it buycut-api-postgres-1 createdb --username=${POSTGRES_USER} ${POSTGRES_DATABASE}

.PHONY: dropdb
dropdb: ## delete the database
	docker exec -it buycut-api-postgres-1 dropdb ${POSTGRES_DATABASE} -U ${POSTGRES_USER} 

.PHONY: injection
injection: ## generate dependency injection code using Wire
	wire gen github.com/ariefro/buycut-api/internal/initializer

.PHONY: run
run: ## run the API server
	APP_ENV=local air
