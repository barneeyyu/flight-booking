.PHONY: all build run seed clean

APP_NAME := flight-booking
SEED_APP_NAME := seed
DB_FILE := flights.db

all: build

build:
	rm -f $(APP_NAME) $(SEED_APP_NAME)
	@echo "Building $(APP_NAME)..."
	go build -o $(APP_NAME) main.go
	@echo "Building $(SEED_APP_NAME)..."
	go build -o $(SEED_APP_NAME) cmd/seed/main.go
	@echo "Build complete."

run:
	@echo "Running $(APP_NAME)..."
	./$(APP_NAME)

seed:
	@echo "Running $(SEED_APP_NAME) to seed data..."
	./$(SEED_APP_NAME)

clean:
	@echo "Cleaning up..."
	rm -f $(APP_NAME) $(SEED_APP_NAME) $(DB_FILE)
	@echo "Cleanup complete."
