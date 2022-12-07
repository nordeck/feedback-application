PKG_LIST := $(shell go list ./... | grep -v /vendor/)

dep:
	@go get -v -d ./...

build: clean dep
	@go build -o out/feedback-api cmd/feedback/main.go

test:
	@go mod tidy
	@DB_HOST=localhost DB_NAME=postgres DB_PASSWORD=postgres DB_TYPE=postgres DB_USER=postgres SSL_MODE=disable DB_PORT=5432 JWT_SECRET=someArbitraryString OIDC_VALIDATION_URL='https://some.url/verify/user' MATRIX_SERVER_NAME='domain.tld' UVS_AUTH_TOKEN=someArbitraryString go test -v ${PKG_LIST}

clean:
	@go clean -testcache
	@rm -rf out

