default: run

minigit:
	go build -o minigit cmd/minigit/main.go

run: minigit
	./minigit

test:
	go test ./... -v

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out
