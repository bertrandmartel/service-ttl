clean:
	@rm -rf dist
	@mkdir -p dist

build: clean
	go build -o service-ttl

run:
	goimports -w .
	gofmt -s -w .
	go run ./main.go

install:
	go get -u ./...

test:
	go test ./... --cover

