gen-normal:
	@echo "Generating normal"
	DEBUG_PGFM=1 protoc \
			-I./normal \
			-I../proto \
			--go_out=paths=source_relative:./normal \
			--fieldmask_out=paths=source_relative,lang=go:./normal \
			./normal/user.proto