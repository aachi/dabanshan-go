NAME=DABANSHAN

default: build

deps: 
	go get -v google/protobuf/timestamp.proto

build: deps