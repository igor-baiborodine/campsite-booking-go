###############################################################################
# Variables
###############################################################################
APP_DOCKER_IMAGE = ibaiborodine/campsite-booking-go

###############################################################################
# Target: image-build
###############################################################################
.PHONY: image-build
image-build:
	@docker build --no-cache -t $(APP_DOCKER_IMAGE) -f ./docker/Dockerfile .
	@docker image ls | grep $(APP_DOCKER_IMAGE)

###############################################################################
# Target: compose-up
###############################################################################
.PHONY: compose-up
compose-up:
	docker compose -f docker/docker-compose.yml -p campsite-booking-go up -d --build

###############################################################################
# Target: compose-up-postgres
###############################################################################
.PHONY: compose-up-postgres
compose-up-postgres:
	docker compose -f docker/docker-compose.yml -p campsite-booking-go up -d postgres

###############################################################################
# Target: compose-down
###############################################################################
.PHONY: compose-down
compose-down:
	docker compose -p campsite-booking-go down -v --remove-orphans
