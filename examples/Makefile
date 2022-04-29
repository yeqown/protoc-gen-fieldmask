gen-pb:
	@echo "Generating pb"
	DEBUG_PGFM=1 protoc \
			-I./pb \
			-I../proto \
			--go_out=paths=source_relative:./pb \
			--fieldmask_out=paths=source_relative,lang=go:./pb \
			./pb/user.proto

gen-common:
	protoc \
    	-I./pb \
    	-I../proto \
    	--go_out=paths=source_relative:./pb \
		./pb/common/common.proto