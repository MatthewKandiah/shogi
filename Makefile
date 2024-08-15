run:
	templ generate
	zig build install
	go run cmd/main.go
