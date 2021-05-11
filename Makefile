build:
	go build -o bin/loadbalancer.out main.go

test:
	go test ./...

run:
	go run main.go
