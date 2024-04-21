.PHONY: install-tools
install-tools:
	echo installing tools
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/vektra/mockery/v2@latest
	echo done

.PHONY: generate
generate:
	echo running code generation
	buf generate
	mockery --quiet --dir ./internal -r --all --inpackage --case underscore
	echo done
