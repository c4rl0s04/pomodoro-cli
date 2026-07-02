.PHONY: build test lint clean run

APP_NAME = pomodoro-cli

build:
	go build -o $(APP_NAME) .

test:
	go test -v ./...

lint:
	golangci-lint run

clean:
	rm -f $(APP_NAME)
	go clean

run: build
	./$(APP_NAME) start
