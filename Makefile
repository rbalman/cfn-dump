build:
	go mod download && go build -o main .

extract-dependencies: build
	STACK_PREFIX=$(STACK_PREFIX) ./main
