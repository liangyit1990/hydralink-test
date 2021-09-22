.PHONY: test build push run db redis setup migrate docker-login

APP_PWD := "$(shell pwd)"
TIMESTAMP := $(shell date +'%s')
GIT_COMMIT := $(shell git rev-parse --short HEAD)

# ----------------------------
# Golang
# ----------------------------
# update-vendor updates the vendor folder
update-vendor:
	go mod vendor
	go mod tidy

# test executes project tests in a golang container
test: dbredo
	go test ./... --coverprofile=coverage.out && \
	go tool cover -html=coverage.out

# runs the application
run:
	export $(shell cat .env.dev) && \
	go build ./cmd/serverd/main.go && ./main

# ----------------------------
# Docker compose
# ----------------------------
COMPOSE := docker-compose -f docker-compose.yml

# compose-down stops and removes all containers and resources associated to docker-compose.yml
compose-down:
	$(COMPOSE) down

compose-run:
	$(COMPOSE) pull go nginx
	$(COMPOSE) up go nginx
# ----------------------------
# Docker
# ----------------------------
# docker-login logs into docker account for private repository
docker-login:
	docker login -u $(DOCKER_USERNAME) -p $(DOCKER_PASSWORD)

push-api: docker-login
	docker tag hydralink-api_go hydralinkapp/hydralink-api:latest && \
	docker push hydralinkapp/hydralink-api:latest

push-web: docker-login
	docker tag hydralink-web hydralinkapp/hydralink-web:latest && \
	docker push hydralinkapp/hydralink-web:latest

# ----------------------------
# Setup and Teardown
# ----------------------------
# setup creates environment
setup: db migrate

# teardown destroys environment
teardown:
	$(COMPOSE) down -v
	$(COMPOSE) rm --force --stop -v

# migrate runs db migration
migrate:
	$(COMPOSE) run --rm -v $(APP_PWD)/data/migrations:/migrations migrate \
	sh -c 'until ./migrate -path /migrations -database $$DATABASE_URL up; do >&2 echo "Postgres is unavailable - sleeping"; sleep 1; done; >&2 echo "Migrate is done..."'

# db starts db container
db:
	$(COMPOSE) up -d db

# dbdrop resets db container
dbdrop:
	$(COMPOSE) run --rm -v $(APP_PWD)/data/migrations:/migrations migrate \
	sh -c './migrate -path /migrations -database $$DATABASE_URL drop'

# dbredo resets db container and performs db migration
dbredo: dbdrop migrate

# ----------------------------
# Binary Build
# ----------------------------
build: build-api build-app

# build executes build in a golang container
build-api:
	@GOOS=linux GOARCH=amd64 go build -v ./cmd/serverd && \
	docker build . -t hydralink-api_go
build-web:
	rm -rf ./nginx/build && \
	cp -r ../web/build ./nginx/build && \
	docker build ./nginx -t hydralink-web

