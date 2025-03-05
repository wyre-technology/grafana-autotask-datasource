.PHONY: build clean test dev

# Build the plugin
build:
	go build -o dist/autotask-datasource ./pkg/main.go

# Clean build artifacts
clean:
	rm -rf dist/

# Run tests
test:
	go test ./...

# Development mode - builds and watches for changes
dev:
	go run pkg/main.go

# Build for different platforms
build-all: clean
	GOOS=linux GOARCH=amd64 go build -o dist/autotask-datasource-linux-amd64 ./pkg/main.go
	GOOS=darwin GOARCH=amd64 go build -o dist/autotask-datasource-darwin-amd64 ./pkg/main.go
	GOOS=windows GOARCH=amd64 go build -o dist/autotask-datasource-windows-amd64.exe ./pkg/main.go

# Install the plugin to Grafana
install: build
	mkdir -p $(GRAFANA_PLUGINS_DIR)/autotask-datasource
	cp -r dist/* $(GRAFANA_PLUGINS_DIR)/autotask-datasource/ 