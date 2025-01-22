.PHONY: run
run:
	go run cmd/main.go

.PHONY: gen
gen:
	go get -u ./...
	go mod tidy
	go generate ./...

.PHONY: test
test:
	go test -v ./...