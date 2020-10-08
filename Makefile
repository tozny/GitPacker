BINARY_NAME=gitpacker
ARCHIVE_DIR=archive

# Mark these targets as not being file based
.PHONY: all lint build clean run version

# Default target executed if no arguments
# provided to make.
all: lint clean build run

# target for linting and formating
# source code and dependency config
lint:
	go fmt ./...
	go vet ./...
	go mod tidy

# target for building binary for the
# local platform architecture from source
build:
	go build -o $(BINARY_NAME)

# target for cleaning up any build artifacts
clean:
	go clean ./...
	rm -rf $(ARCHIVE_DIR)/
	rm $(BINARY_NAME) || true
	rm packedByGitPacker.zip || true

# target for running program based off
# latest source code
run:
	go run main.go

# target for tagging and publishing a
# new version of the program
# run like make version=X.Y.Z
version:
	git tag v${version}
	git push origin v${version}
