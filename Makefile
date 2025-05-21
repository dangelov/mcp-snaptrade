build:
	go build -o bin/cli cmd/cli/main.go
	chmod +x bin/cli

manage:
	cd cmd/manage && go run main.go

test:
	cd cmd/data && go run main.go
