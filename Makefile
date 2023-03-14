GO=${GOROOT}/bin/go

swagger-init:
	${GO} get github.com/swaggo/swag/cmd/swag
	${GO} get github.com/swaggo/http-swagger
	${GO} get github.com/alecthomas/template
	swag init -g cmd/api/api-service.go
	${GO} mod tidy -compat=1.17
swagger:
	swag init -g ./cmd/api/api-service.go
git-update:
	git rm -rf --cached .
	git add .
test:
	${GO} test -race -short -vet=off ./...
build:
	#CGO_ENABLED=0 GOARCH=amd64 GOOS=windows ${GO} build -ldflags '-extldflags "-static"' -o dist/win cmd/api/api-service.go
	#CGO_ENABLED=0 GOARCH=amd64 GOOS=linux ${GO} build -ldflags '-extldflags "-static"' -o dist/lin cmd/api/api-service.go
	CGO_ENABLED=0 GOARCH=arm64 GOOS=darwin ${GO} build -tags kqueue -ldflags '-extldflags "-static"' -o dist/mac cmd/api/api-service.go
test-zip:
	mkdir -p zipdir
	zip -r -Z bzip2 zipdir/test.zip docs

mongo:
	docker run --name mongo -e MONGODB_DATABASE=db -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=qwerty -p 27017:27017 -v /Users/eminshakh/data/mongodb-data:/data/db -d arm64v8/mongo


run-docker:
	docker run -d \
  	-it \
  	--name vimb-loader \
  	-v $(pwd)/config.json:/config.json \
  	-v $(pwd)/db:/db \
  	-v $(pwd)/logs:/logs \
  	-v $(pwd)/s3-buckets:/s3-buckets \
  	-p 8000:8000 \
  	ghcr.io/advancemg/vimb-loader:1.0.0