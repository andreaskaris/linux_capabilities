.PHONY: ambient
ambient:
	mkdir -p bin
	go build -o bin/ambient cmd/ambient/main.go

.PHONY: http
http:
	mkdir -p bin
	go build -o bin/http cmd/http/main.go

.PHONY: build
build: ambient http

.PHONY: set-cap-example
set-cap-example:
	sudo setcap "cap_net_bind_service+p cap_chown+p" bin/ambient
