PKG = github.com/dgo-/gopml


all: install



install:
	go install ${PKG}


build:
	go build ${PKG}


test:
	go test -race -coverprofile=c.out -covermode=atomic

cover: test
	go tool cover -html=c.out -o coverage.html

