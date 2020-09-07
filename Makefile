.PHONY: build start stop logs down

build: ## Build docker containers
	docker-compose build

start: ## Start docker containers
	docker-compose up -d

stop: ## Stop docker containers
	docker-compose stop

app-logs: ## Show app logs
	docker-compose logs -f app

down: ## Down docker containers
	docker-compose down --volumes