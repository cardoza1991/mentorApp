# Build variables
APP_NAME := mentorapp
LAUNCHER_NAME := $(APP_NAME)_launcher
BUILD_DIR := dist
POSTGRES_DIR := $(BUILD_DIR)/pgsql

.PHONY: all clean build package install-postgres

all: package

clean:
	rm -rf $(BUILD_DIR)

build:
	# Create build directory
	mkdir -p $(BUILD_DIR)
	mkdir -p $(BUILD_DIR)/conf
	
	# Build the main application
	go build -tags 'postgres' -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server
	
	# Build the launcher
	go build -o $(BUILD_DIR)/$(LAUNCHER_NAME) ./cmd/launcher

install-postgres:
	# Create postgres directory
	mkdir -p $(POSTGRES_DIR)/bin
	mkdir -p $(POSTGRES_DIR)/lib
	
	# Copy PostgreSQL binaries and libraries
	cp /usr/lib/postgresql/15/bin/postgres $(POSTGRES_DIR)/bin/
	cp /usr/lib/postgresql/15/bin/pg_ctl $(POSTGRES_DIR)/bin/
	cp /usr/lib/postgresql/15/bin/initdb $(POSTGRES_DIR)/bin/
	cp /usr/lib/postgresql/15/bin/pg_isready $(POSTGRES_DIR)/bin/
	cp /usr/lib/postgresql/15/bin/psql $(POSTGRES_DIR)/bin/
	cp -r /usr/lib/postgresql/15/lib/* $(POSTGRES_DIR)/lib/
	
	# Set correct permissions
	chmod +x $(POSTGRES_DIR)/bin/*

setup-postgres:
	# Install PostgreSQL if not already installed
	if ! dpkg -l | grep -q postgresql-15; then \
		echo "Installing PostgreSQL 15..." && \
		sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $$(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list' && \
		wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add - && \
		sudo apt-get update && \
		sudo apt-get install -y postgresql-15; \
	fi

package: clean build setup-postgres install-postgres
	# Copy configuration files
	mkdir -p $(BUILD_DIR)/conf
	cp internal/config/app.conf $(BUILD_DIR)/conf/ || touch $(BUILD_DIR)/conf/app.conf
	
	# Copy necessary files
	cp -r static $(BUILD_DIR)/
	cp -r views $(BUILD_DIR)/
	cp -r migrations $(BUILD_DIR)/
	
	# Copy README and license files
	cp README.md $(BUILD_DIR)/
	
	# Create final archive
	tar -czf $(APP_NAME)-package.tar.gz -C $(BUILD_DIR) .

.PHONY: run
run:
	./$(BUILD_DIR)/$(LAUNCHER_NAME)