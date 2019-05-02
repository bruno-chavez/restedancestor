all: test hard-lint install

test:
	go test ./...

hard-lint:
	gometalinter --enable-all -D lll -t --sort=severity ./...

install:
	go install ./...
