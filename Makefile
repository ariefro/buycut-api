.PHONY: injection
injection: ## generate dependency injection code using Wire
	wire gen github.com/ariefro/buycut-api/internal/initializer

.PHONY: run
run: ## run the API server
	APP_ENV=local air
