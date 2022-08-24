.PHONY: apiserver test build clean

apiserver:
	echo "building apiserver"; \
	go build -o apiserver main.go

build: apiserver

test:
	go test -v; \
	cd dbactions && go test -v

clean:
	rm apiserver; \
	rm api_data.db; \
	rm dbactions/dbactions_test.db
