build:
	go mod download && go build -o main .

test:
	go test ./...
