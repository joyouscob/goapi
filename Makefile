dev:
	@reflex -r '.go' -s -- go run main.go

run:
	go get github.com/cespare/reflex
	reflex -r '\.go$$' -s go run .