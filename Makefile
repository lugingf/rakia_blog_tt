.PHONY: lint, test, check

test:
	go test -v ./... 

test-n:
	source .env.local && go test -v -count=1 -p=1 ./... -run $(name)

lint:
	golangci-lint run -v --skip-dirs=test

check: lint test

run:
	make build && source .env.local && ./event_meld

build:
	go build -ldflags "-s -w" -o blog_tt;

.PHONY: models
models:
	swagger generate model -f ./api/api.json -t ./handler
