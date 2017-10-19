DEST = ./
PROTOBUFS = $(wildcard *.proto)
PROTOC_FLAGS  = -I .

all:
	protoc \
	--proto_path=./ \
	--go_out=plugins=grpc:./ ./*.proto