.PHONY:

build-ex:
	cd ./extension && GOOS=linux GOARCH=amd64 go build -o ../bin/extensions/telemetry-api-extension main.go
	chmod +x bin/extensions/telemetry-api-extension
	cd bin && zip -r extension.zip extensions/

build-func:
	cd ./function && zip ../function.zip main.py

role:
	aws iam create-role --role-name telemetry-api-function-role --assume-role-policy-document '{"Version": "2012-10-17","Statement": [{ "Effect": "Allow", "Principal": {"Service": "lambda.amazonaws.com"}, "Action": "sts:AssumeRole"}]}' --query 'Role.Arn' --output text
	aws iam attach-role-policy --role-name telemetry-api-function-role --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole

