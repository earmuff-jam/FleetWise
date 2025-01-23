
include .env

# load swagger for api contract
# for some weird reason swagger is not working via makefile
# navigate to apilayer and run cmd - swagger generate spec -o swagger.yaml --scan-models
swagger:
	cd apilayer && \
	swagger generate spec -o swagger.yaml --scan-models

# load api in development instance
startserver:
	cd apilayer && \
	GO111MODULE=on \
	go run api.go

# api unit tests without coverage printed out
gotest:
	./main.sh -T

# api unit test that can be run anytime during test period. must run against test db not main
apitest:
	cd apilayer/handler && \
	go test -coverprofile=../logs/coverage.out ./... && go tool cover -func=../logs/coverage.out

# load ui for the first time; fresh install and start ui client
runClient:
	cd frontend && \
	yarn && \
	yarn dev

# load ui in development instance
startclient:
	cd frontend && \
	yarn dev

# run the lint in frontend code
linter:
	cd frontend && \
	yarn lint