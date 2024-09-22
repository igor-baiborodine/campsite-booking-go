###############################################################################
# Variables
###############################################################################
APP_DOCKER_IMAGE = ibaiborodine/campsite-booking-go

###############################################################################
# Target: image-build-app
###############################################################################
.PHONY: image-build-app
image-build-app:
	@docker build --no-cache -t $(APP_DOCKER_IMAGE) -f ./docker/Dockerfile .
	@docker image ls | grep $(APP_DOCKER_IMAGE)

###############################################################################
# Target: compose-up-all
###############################################################################
.PHONY: compose-up-all
compose-up-all:
	docker compose -f docker/docker-compose.yml -p campsite-booking-go up -d --build

###############################################################################
# Target: compose-up-postgres
###############################################################################
.PHONY: compose-up-postgres
compose-up-postgres:
	docker compose -f docker/docker-compose.yml -p campsite-booking-go up -d postgres

###############################################################################
# Target: compose-down-all
###############################################################################
.PHONY: compose-down-all
compose-down-all:
	docker compose -p campsite-booking-go down -v --remove-orphans
