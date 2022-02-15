go build
go build --buildmode=plugin ./plugin_beta/plugins/cq-chatlogger/
cp ./cq-chatlogger.so ~/.config/fastbuilder/plugins_beta/
echo amd64,ok

