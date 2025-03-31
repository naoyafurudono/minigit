default: run

minigit:
	go build -o minigit cmd/minigit/main.go

run: minigit
	./minigit
