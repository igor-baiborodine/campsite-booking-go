.PHONY: install-tools
install-tools:
	echo installing tools
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4.0
	go install github.com/bufbuild/buf/cmd/buf@v1.34.0
	go install github.com/vektra/mockery/v2@v2.43.2
	echo done

.PHONY: generate
generate:
	echo running code generation
	go generate
	echo done
