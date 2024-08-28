-include k8s/local-k8s.mk

################################################################################
# Variables                                                                    #
################################################################################
PROTOC_GEN_GO_VERSION = v1.34.2
PROTOC_GEN_GO_GRPC_VERSION = v1.4.0
MOCKERY_VERSION = v2.43.2
GOIMPORTS_VERSION = v0.22.0
GOLINES_VERSION = v0.12.2
GOFUMPT_VERSION = v0.6.0
GOLANGCI_LINT_VERSION = v1.60.3

################################################################################
# Target: init-proto
################################################################################
.PHONY: init-proto
init-proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)

################################################################################
# Target: init-mock
################################################################################
.PHONY: init-mock
init-mock:
	go install github.com/vektra/mockery/v2@$(MOCKERY_VERSION)

################################################################################
# Target: init-format
################################################################################
.PHONY: init-format
init-format:
	go install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)
	go install github.com/segmentio/golines@$(GOLINES_VERSION)
	go install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)

################################################################################
# Target: init-golangci-lint
################################################################################
.PHONY: init-golangci-lint
init-golangci-lint:
	sudo curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)
	golangci-lint --version

################################################################################
# Target: install-tolls
################################################################################
.PHONY: install-tools
install-tools: init-proto init-mock init-format init-golangci-lint

################################################################################
# Target: mod-tidy
################################################################################
.PHONY: mod-tidy
mod-tidy:
	go mod tidy

################################################################################
# Target: gen-proto
################################################################################
.PHONY: gen-proto
gen-proto:
	buf generate

################################################################################
# Target: gen-mock
################################################################################
.PHONY: gen-mock
gen-mock:
	mockery --quiet --dir ./internal -r --all --inpackage --case underscore

################################################################################
# Target: format
################################################################################
.PHONY: format
format:
	golines -w --ignore-generated . && gofumpt -w .

################################################################################
# Target: format-proto
################################################################################
.PHONY: format-proto
format-proto:
	buf format -w

################################################################################
# Target: lint
################################################################################
.PHONY: lint
lint:
	golangci-lint run

################################################################################
# Target: lint-proto
################################################################################
.PHONY: lint-proto
lint-proto:
	buf lint

################################################################################
# Target: test
################################################################################
.PHONY: test
test:
	go test -race $(COVERAGE_OPTS) ./internal/...

################################################################################
# Target: test-integration
################################################################################
.PHONY: test-integration
test-integration:
	go test -tags=integration ./internal/...

################################################################################
# Target: check
################################################################################
.PHONY: check
check: mod-tidy format test test-integration lint

################################################################################
# Target: check-mod-diff                                                       #
################################################################################
.PHONY: check-mod-diff
check-mod-diff:
	git diff --exit-code ./go.mod
	git diff --exit-code ./go.sum

################################################################################
# Target: check-proto-diff
################################################################################
.PHONY: check-proto-diff
check-proto-diff:
	git diff --exit-code ./campgroundspb # generated pb

################################################################################
# Target: check-mock-diff
################################################################################
.PHONY: check-mock-diff
check-mock-diff: check-format-diff

################################################################################
# Target: check-format-diff
################################################################################
.PHONY: check-format-diff
check-format-diff:
	git diff --exit-code . ':!campgroundspb' # not generated pb
