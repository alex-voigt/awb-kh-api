help:
	@echo "Available commands:"
	@echo "	run                - runs the api"
	@echo "	watch              - runs the api with hot reload"
	@echo "	test               - runs the tests"
	@echo "	docker-build       - builds the docker container"
	@echo "	docker-run         - runs the docker container"
	@echo ""

.PHONY: run
run:
	go build && ./awb-kh-api

.PHONY: watch
watch:
	go get -u github.com/cosmtrek/air
	air

.PHONY: test
test:
	go test ./... -timeout 30s -v -cover

.PHONY: docker-build
docker-build:
	docker build -t awb-kh-api .

.PHONY: docker-run
docker-run:
	 docker run -p 8010:8010 --name awb-kh-api --rm -it awb-kh-api:latest