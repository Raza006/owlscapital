# Common tasks for running the project with Docker Compose

COMPOSE_FILE := services/docker-compose.yml
COMPOSE := docker compose -f $(COMPOSE_FILE)

DISCORD_SERVICE := discordbot
POSTGRES_SERVICE := postgres

.DEFAULT_GOAL := help

.PHONY: help up down up-postgres restart restart-postgres logs logs-postgres build ps pull clean recreate

help:
	@echo "Common tasks:"
	@echo "  make up         Build and start containers (detached)"
	@echo "  make down       Stop and remove containers"
	@echo "  make up-postgres  Start only the postgres service"
	@echo "  make restart    Restart only the discordbot service"
	@echo "  make restart-postgres  Restart only the postgres service"
	@echo "  make logs       Tail logs of the discordbot service"
	@echo "  make logs-postgres   Tail logs of the postgres service"
	@echo "  make build      Build the discordbot image"
	@echo "  make ps         Show compose services status"
	@echo "  make pull       Pull remote images (if any)"
	@echo "  make clean      Stop and remove containers, networks and volumes"
	@echo "  make recreate   Recreate containers (force)"

up:
	$(COMPOSE) up -d --build

down:
	$(COMPOSE) down

up-postgres:
	$(COMPOSE) up -d $(POSTGRES_SERVICE)

restart:
	$(COMPOSE) restart $(DISCORD_SERVICE)

restart-postgres:
	$(COMPOSE) restart $(POSTGRES_SERVICE)

logs:
	$(COMPOSE) logs -f $(DISCORD_SERVICE)

logs-postgres:
	$(COMPOSE) logs -f $(POSTGRES_SERVICE)

build:
	$(COMPOSE) build $(DISCORD_SERVICE)

ps:
	$(COMPOSE) ps

pull:
	$(COMPOSE) pull

clean:
	$(COMPOSE) down -v

recreate:
	$(COMPOSE) up -d --force-recreate


