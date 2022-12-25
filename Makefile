PORT=9990
ADDR ?= wow-server:${PORT}

lint:
	go vet -v ./...

test: lint
	go test -race -v ./...

docker-build-srv:
	docker build -t wow-server \
		--build-arg addr=${ADDR} \
		--build-arg path=server \
		.

docker-build-client:
	docker build -t wow-client \
		--build-arg addr=${ADDR} \
		--build-arg path=client \
		.

docker-network:
	docker network create wow || exit 0

start-srv: docker-network
	docker run \
		-p ${PORT}:${PORT} \
		--rm -d \
		--cpus 1 \
		--memory 300M \
		--name wow-server \
		--network wow \
		wow-server

start-client: docker-network
	docker run \
		--rm \
		--network wow \
		wow-client

docker-clean:
	docker stop wow-server wow-client || exit 0
	docker rmi -f wow-client wow-server || exit 0
	docker network rm wow || exit 0

docker-build: docker-clean docker-network docker-build-srv docker-build-client

docker-run: start-srv start-client

demo: docker-network docker-build-srv docker-build-client start-srv start-client docker-clean