FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp

## upd: starts all containers in the background without forcing build
upd:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
upbd:
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_front: builds the frone end binary
fe.build:
	@echo "Building front end binary..."
	cd ./front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

## start: starts the front end
fe.start:
	@if [ ! -f "./front-end/${FRONT_END_BINARY}" ]; then \
        make build_front; \
    fi
	@echo "Starting front end"
	cd ./front-end && ./${FRONT_END_BINARY} &

## stop: stop the front end
fe.stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"

# install & attach delve debugger to authentication service
debug.auth:
	docker exec -i svc-auth /bin/sh < ./authentication-service/delve.sh &

# install & attach delve debugger to broker service
debug.broker:
	docker exec -i svc-broker /bin/sh < ./broker-service/delve.sh &
