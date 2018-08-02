all: test build

run: build start

init_db:
	@echo " >> initialize database"
	@go run cmd/dbTest/main.go

test:
	@echo " >> running tests"
	@go test -v -race ./...

build:
	@echo " >> building binaries"
	@go build
	
start:
	@echo " >> starting binaries"
	@./oms-cart
