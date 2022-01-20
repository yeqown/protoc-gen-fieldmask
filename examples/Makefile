gen-normal:
	@echo "Generating normal"
	protoc \
			-I./normal \
			-I../proto \
			--go_out=paths=source_relative:./normal \
			--fieldmask_out=paths=source_relative:./normal \
			./normal/user.proto