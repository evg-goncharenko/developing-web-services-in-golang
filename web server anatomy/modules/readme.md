go mod init myapp
# go mod init github.com/rvasily/myapp
go build
go mod download
go mod verify
go mod tidy

go build  -o ./bin/myapp ./cmd/myapp
go test -v -coverpkg=./... ./...

go mod vendor
go build -mod=vendor -o ./bin/myapp ./cmd/myapp
go test -v -mod=vendor -coverpkg=./... ./...