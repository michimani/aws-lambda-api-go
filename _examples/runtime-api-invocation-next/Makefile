.PHONY: syncmod
syncmod:
	go mod tidy
	go mod vendor

.PHONY: build run build-run
build: syncmod
	docker build -t runtime-api-invocation-next:local .

run:
	docker run \
	--rm \
	-p 9000:8080 \
	runtime-api-invocation-next:local

build-run: build run