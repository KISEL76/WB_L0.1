COMPOSE = docker-compose
PRODUCER_PATH=cmd/producer/

all: launch produce

produce:
	@cd $(PRODUCER_PATH) && go run main.go
	@echo "Данные отправлены в брокер!"

build:
	$(COMPOSE) build --no-cache

launch:
	@$(COMPOSE) up -d 

logs:
	@$(COMPOSE) logs --tail=100

stop:
	@$(COMPOSE) down

stop-v:
	@$(COMPOSE) down -v 
	



