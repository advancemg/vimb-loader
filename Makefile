GO=${GOROOT}/bin/go

swagger-init:
	${GO} get github.com/swaggo/swag/cmd/swag
	${GO} get github.com/swaggo/http-swagger
	${GO} get github.com/alecthomas/template
	swag init -g cmd/api/api-service.go
	go mod tidy
git-update:
	git rm -rf --cached .
	git add .
build:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows ${GO} build -ldflags '-extldflags "-static"' -o dist/win cmd/api/api-service.go
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux ${GO} build -ldflags '-extldflags "-static"' -o dist/lin cmd/api/api-service.go
