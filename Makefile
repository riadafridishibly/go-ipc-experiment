.PHONY: makecert
makecert:
	bash makecert.sh test@example.local

.PHONY: build
build-counter:
	go build -o counter/counter counter/main.go

.PHONY: run-server makecert
run-server: build-counter
	go run -race server/main.go

.PHONY: run-client
run-client:
	go run -race client/main.go