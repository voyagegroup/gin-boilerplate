NPM:=$(shell which npm)
NODE_MODULES := node_modules/.bin
WEBPACK:=$(NODE_MODULES)/webpack
ESLINT:=$(NODE_MODULES)/eslint
WEBPACK_JS:=client/webpack/webpack.config.js
DB_CONTAINER_NAME:=treasure-mysql
DBNAME:=treasure
ENV:=development

deps:
	which godep || go get github.com/tools/godep
	godep restore
	go get github.com/rubenv/sql-migrate/...
	npm install

test:
	godep go test -v ./...

integration-test:
	godep go test -v ./... -tags=integration

save:
	godep save

build:
	$(WEBPACK) -p --config $(WEBPACK_JS)

watch:
	$(WEBPACK) -w --config $(WEBPACK_JS)

lint:
	$(ESLINT) client/**/*.jsx

fix:
	$(ESLINT) --fix client/**/*.jsx

migrate/init:
	mysql -u root -h localhost --protocol tcp -e "create database \`$(DBNAME)\`" -p

migrate/up:
	sql-migrate up -env=$(ENV)

migrate/down:
	sql-migrate down -env=$(ENV)

migrate/status:
	sql-migrate status -env=$(ENV)

migrate/dry:
	sql-migrate up -dryrun -env=$(ENV)

docker/build: Dockerfile docker-compose.yml
	docker-compose build

docker/start:
	docker-compose up -d

docker/logs:
	docker-compose logs

docker/clean:
	docker-compose rm

docker/bash:
	docker exec -it $(shell docker-compose ps -q) bash
