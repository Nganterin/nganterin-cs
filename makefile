ARGS=$(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

%:
	@:

# Copy blueprint to new-feature
cp-blueprint:
	cp -r api/blueprint api/$(ARGS)
	find api/$(ARGS) -type f -exec sed -i "s/blueprint/$(ARGS)/g" {} +
	find api/$(ARGS) -type f -name 'blueprint_*' -exec bash -c 'mv "$$0" "$${0%/*}/$(ARGS)_$${0##*/blueprint_}"' {} \;

# Run the application
run:
	go mod tidy
	go run cmd/nganterin-cs/main.go

# Build the application binary
build:
	go build -o bin/nganterin-cs ./cmd/nganterin-cs

# Run migrations
migrate:
	go run cmd/migrate/migrate.go

# Clean the build (remove binaries and build artifacts)
clean:
	rm -f bin/nganterin-cs

# Change environment to production
env-prod:
	rm .env
	cp .env.prod .env

# Change environment to development
env-dev:
	rm .env
	cp .env.dev .env

# Generate wire dependencies
wire:
	wire gen ./injectors

# Generate wire dependencies for internal
wire-internal:
	wire gen ./internal/injectors