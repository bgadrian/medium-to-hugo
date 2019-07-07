# Makefile
source := main.go

export GO111MODULE = on

pre:
	mkdir -p ./build/
	go get -d ./

build: pre
	go build -o ./build/binary $(source)

buildall: pre
	mkdir -p ./build/mediumtohugo/windows
	mkdir -p ./build/mediumtohugo/linux
	mkdir -p ./build/mediumtohugo/macos
	GOOS=darwin GOARCH=amd64 go build -o ./build/mediumtohugo/macos/mediumtohugo $(source)
	GOOS=linux GOARCH=amd64 go build -o ./build/mediumtohugo/linux/mediumtohugo $(source)
	GOOS=windows GOARCH=amd64 go build -o  ./build/mediumtohugo/windows/mediumtohugo.exe $(source)
	cd ./build && tar -czf ./mediumtohugo.tar.gz ./mediumtohugo/
	@echo "publish to gihub: $ hub release create -a ./build/mediumtohugo.tar.gz -m 'v0.X' v0.X"
