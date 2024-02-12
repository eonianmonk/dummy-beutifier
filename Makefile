build:
	go build -o beautyf cmd/main.go

serve:
	./beautyf run

tests:
	go test ./test