export CC=x86_64-w64-mingw32-gcc
export CGO_ENABLED=1
export GOARCH=amd64
go build
go build --buildmode=plugin ./plugin_beta/plugins/cq-chatlogger/
cp ./cq-chatlogger.so ~/.config/fastbuilder/plugins_beta/
echo amd64,ok
export GOARM=7
export GOARCH=arm64
export GOOS=linux
go build
