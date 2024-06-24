.PHONY: install-tools
install-tools:
	echo installing tools
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4.0
	go install github.com/bufbuild/buf/cmd/buf@v1.34.0
	go install github.com/vektra/mockery/v2@v2.43.2
	echo done

.PHONY: format
format:
	gofmt -s -w .

.PHONY: format-diff
format-diff:
	gofmt -s -w . && git diff --exit-code

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: tidy-diff
tidy-diff:
	go mod tidy && git diff --exit-code

.PHONY: download
download:
	go mod download

.PHONY: verify
verify:
	go mod verify

.PHONY: generate
generate:
	go generate ./...

.PHONY: generate-diff
generate-diff:
	go generate ./... && git diff --exit-code

.PHONY: build
build:
	go build -v ./...

.PHONY: unit-tests
unit-tests:
	go test -v ./...

.PHONY: integration-tests
integration-tests:
	go test ./... -tags integration
