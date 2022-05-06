gen-issue17:
	DEBUG_PGFM=1 protoc \
			-I. \
			-I../../proto \
			--go_out=paths=source_relative:. \
			--fieldmask_out=paths=source_relative,lang=go:. \
			./issue17.proto