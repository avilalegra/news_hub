test/unit:
	@GOFLAGS="-count=1" go test -short ./...

test/int:
	@GOFLAGS="-count=1" go test -run Integration -p 1 ./...

test/all: test/unit test/int