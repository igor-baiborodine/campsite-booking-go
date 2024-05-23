.PHONY: install-tools
install-tools:
	echo installing tools
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/vektra/mockery/v2@latest
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	echo done

.PHONY: generate
generate:
	echo running code generation
	go generate
	echo done
