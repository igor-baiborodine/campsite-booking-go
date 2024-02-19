.PHONY: install-tools
install-tools:
	echo installing tools
	go install \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc
	echo done

.PHONY: generate
generate:
	echo running code generation
	buf generate
	echo done
