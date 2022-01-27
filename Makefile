debug_path=$(shell which protoc-gen-debug)

install:
	go install ./

test:
	@echo "test"
	go test -v ./... --count=1

gen-fm-pb:
	protoc \
		-I=./proto \
		--go_out=paths=source_relative:./proto \
		./proto/fieldmask/option.proto

prepare-debug:
	- mkdir internal/module/debugdata
	protoc \
		-I=./examples/pb \
		-I=./proto \
		--plugin=protoc-gen-debug=${debug_path} \
		--debug_out="./internal/module/debugdata:." \
		./examples/pb/user.proto