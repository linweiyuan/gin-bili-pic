CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" main.go

docker build -t linweiyuan/gin-bili-pic .
