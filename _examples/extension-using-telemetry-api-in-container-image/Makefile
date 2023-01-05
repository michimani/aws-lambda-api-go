.PHONY:

build-ex:
	cd ./extension && GOOS=linux GOARCH=amd64 go build -o ../function/bin/extensions/telemetry-api-extension main.go
	chmod +x function/bin/extensions/telemetry-api-extension
	cd function/bin && zip -r extension.zip extensions/

build-func:
	cd ./function && docker build -t telemetry-api-ex-func:local .

run:
	docker run \
	--rm \
	-p 9000:8080 \
	telemetry-api-ex-func:local