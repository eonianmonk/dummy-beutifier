build:
	go build -o beautyf cmd/main.go

serve:
	./beautyf run

test:
	go test ./test