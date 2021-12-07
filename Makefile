build:
	@docker run --rm -v ${PWD}:/go/src/github.com/authorizer -w '/go/src/github.com/authorizer' golang:1.16 go build -o ./tmp/authorizer ./cmd/cli

run:
	@./tmp/authorizer

test-coverage:
	@go test -race -p=1 -coverprofile ./cover.out ./... && go tool cover -html=./cover.out

test-benchmark:
	- go test ./... -bench=.

fmt:
	- go fmt ./...

scan-cyclomatic-complexity:
	-  gocyclo -over 5 .

go-sec:
	- gosec -out report.json ./...