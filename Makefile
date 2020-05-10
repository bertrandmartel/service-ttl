clean:
	@rm -rf dist
	@mkdir -p dist

build: clean
	go build -o service-ttl
	env GOOS=linux GOARCH=arm GOARM=5 go build -o service-ttl-arm

run:
	goimports -w .
	gofmt -s -w .
	go run ./main.go

install:
	go get -u ./...

test:
	go test ./... --cover

