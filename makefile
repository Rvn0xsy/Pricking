build:
	@echo "Building Pricking...."
	@go build

test:
	@./Pricking -url https://payloads.online
clean:
	@rm -rf ./Pricking
help:
	@echo make build
	@echo make test
	@echo make test
