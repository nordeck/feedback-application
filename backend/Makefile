PKG_LIST := $(shell go list ./... | grep -v /vendor/)

dep:
	@go get -v -d ./...

build: clean dep
	@go build -o out/feedback-api cmd/feedback/main.go

test:
	@go mod tidy
	@DB_HOST=localhost DB_NAME=postgres DB_PASSWORD=postgres DB_TYPE=postgres DB_USER=postgres SSL_MODE=disable DB_PORT=5432 JWT_SIGNATURE=someArbitraryString OIDC_VALIDATION_URL='https://github.com/verify/user' go test -v ${PKG_LIST}

clean:
	@go clean -testcache
	@rm -rf out

