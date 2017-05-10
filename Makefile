BINARY=solution

build:
	go build -o ${BINARY}

install:
	go install
	go get github.com/sirupsen/logrus

clean:
	go clean
	rm -rf solution

.PHONY: build clean install
