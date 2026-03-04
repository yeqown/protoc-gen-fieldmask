.PHONY: genpb

genpb:
	protoc \
		-I . \
		-I ../third_party \
		--go_out=paths=source_relative:. \
		--fieldmask_out=paths=source_relative,lang=go:. \
		./example.proto
