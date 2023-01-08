precondition:
	# you have installed the protoc-gen-fieldmask
	# go install github.com/yeqown/protoc-gen-fieldmask@latest

	# check protoc-gen-fieldmask has been installed, if not, install it
	$(shell if ! command -v protoc-gen-fieldmask &> /dev/null; then \
		echo "protoc-gen-fieldmask not found, install it"; \
		go install github.com/yeqown/protoc-gen-fieldmask@latest; \
	fi)

run: precondition
	DEBUG_PGFM=1 protoc \
			-I. \
			-I../../../proto \
			--go_out=paths=source_relative:. \
			--fieldmask_out=paths=source_relative,lang=go:. \
			./issue21_0.proto ./issue21_1.proto