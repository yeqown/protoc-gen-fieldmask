install:
	go install ./

test:
	@echo "test"
	go test -v ./... --count=1