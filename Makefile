.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/eventbridge/plans handlers/aws/eventbridge/plans/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/eventbridge/subscriptions handlers/aws/eventbridge/subscriptions/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/eventbridge/validate/challenge handlers/aws/eventbridge/validate/challenge/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/eventbridge/new_subscription_email handlers/aws/eventbridge/new_subscription_email/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/eventbridge/create_auth0_user handlers/aws/eventbridge/create_auth0_user/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/apigateway/e2e_webhook handlers/aws/apigateway/e2e_webhook/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/apigateway/subscriptions handlers/aws/apigateway/subscriptions/main.go

test: build
	go test ./...

clean:
	rm -rf ./bin

deploy: build
	sls deploy --force --verbose
