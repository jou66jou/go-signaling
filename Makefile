# 檔名後綴
VERSION := $(shell git describe --tags --always)
OS := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))
TAG := $(VERSION)_$(ARCH)
# 編譯目標
MAINPATH := ./cmd/p2p/
# 輸出檔名
EXEC_FILENAME := /p2p
# 輸出目標
EXEC_PATH := ./exec

# 本地編譯後執行
run : 
	go run $(MAINPATH)/main.go 
build:
	go build -o $(EXEC_PATH)$(EXEC_FILENAME)_$(TAG)_$(OS) $(MAINPATH) 
build.run:
	go build -o $(EXEC_PATH)$(EXEC_FILENAME)_$(TAG)_$(OS) $(MAINPATH) 
	$(EXEC_PATH)$(EXEC_FILENAME)_$(TAG)_$(OS)
delete.all:
	rm $(EXEC_PATH)/*
	
# 給 mac 用
mac.build.linux: 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(EXEC_PATH)$(EXEC_FILENAME)_$(TAG)_linux $(MAINPATH) 
mac.build.exe:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(EXEC_PATH)$(EXEC_FILENAME)_$(TAG)_win.exe $(MAINPATH) 
mac.build.all: build mac.build.linux mac.build.exe

# 給 linux 用
linux.build.mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(EXEC_PATH)$(EXEC_FILENAME)_$(TAG)_linux $(MAINPATH) 
linux.build.exe:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(EXEC_PATH)$(EXEC_FILENAME)_$(TAG)_win.exe $(MAINPATH) 

# 給 windows 用
win.build.mac:
	SET CGO_ENABLED=0
	SET GOOS=darwin
	SET GOARCH=amd64
	go build -o $(EXEC_PATH)$(EXEC_FILENAME)_$(TAG)_darwin $(MAINPATH) 
win.build.linux:
	SET CGO_ENABLED=0
	SET GOOS=linux
	SET GOARCH=amd64
	go build -o $(EXEC_PATH)$(EXEC_FILENAME)_$(TAG)_linux $(MAINPATH) 