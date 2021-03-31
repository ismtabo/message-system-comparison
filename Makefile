define help
Usage: make <command>
Commands:
  help:              Show this help information
  clean:             Clean the project (remove build directory, clean golang packages and tidy go.mod file)
  build-bin:         Build the application into build/bin directory
  build:             Build the application. Orchestrates: build-bin, build-config, build-test and build-chown
  run:               Launch the application
  develenv-up:       Launch the development environment with a docker-compose of the service
  develenv-sh:       Access to a shell of the develenv service.
  develenv-down:     Stop the development environment
endef
export help

.PHONY: help
help:
	@echo "$$help"

.PHONY: clean
clean:
	$(info) 'Cleaning the project'
	rm -rf build/
	go clean
	go mod tidy

.PHONY: build-bin
build-bin:
	$(info) 'Building version: $(BUILD_VERSION)'
	mkdir -p build/bin
	go build -v -ldflags='$(LDFLAGS)' -o build/bin ./...

.PHONY: build
build: build-bin

.PHONY: run
run:
	$(info) 'Launching the service'
	build/bin/redirector

.PHONY: develenv-up
develenv-up:
	$(info) 'Launching the development environment: $(PRODUCT_VERSION)-$(PRODUCT_REVISION)'
	$(DOCKER_COMPOSE) up --build -d

.PHONY: develenv-sh
develenv-sh:
	# $(DOCKER_COMPOSE) exec develenv bash

.PHONY: develenv-down
develenv-down:
	$(info) 'Shutting down the development environment'
	$(DOCKER_COMPOSE) down --remove-orphans

# Functions
info := @printf '\033[32;01m%s\033[0m\n'
get_packages := $$(go list ./... | grep -v test/acceptance)