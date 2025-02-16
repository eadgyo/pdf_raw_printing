
.PHONY: build deploy clean route1 route2 route1local route2local

build:
	go build -o "./build/" ./cmd/...

clean:
	rm -rf build/*
