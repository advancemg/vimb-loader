GO=${GOROOT}/bin/go

swagger-init:
	${GO} get github.com/swaggo/swag/cmd/swag
	${GO} get github.com/swaggo/http-swagger
	${GO} get github.com/alecthomas/template
	swag init -g cmd/api/api-service.go
	go mod tidy
