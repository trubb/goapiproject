.PHONY: apiserver test build clean

apiserver:
	echo "building apiserver"; \
	go build -o apiserver main.go

build: apiserver

test:
	go test -v; \
	rm *_test.db

clean:
	rm apiserver; \
	rm *.db
