run:
	go run ./cmd/habitum

tidy:
	@echo "Formatting .go files..."
	go fmt ./...
	@echo "Tidying module dependencies..."
	go mod tidy
	@echo "Verifying dependencies..."
	go mod verify
