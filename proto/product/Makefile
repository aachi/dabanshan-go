DEST = ./
PROTOBUFS = $(wildcard *.proto)
PROTOC_FLAGS  = -I .

all:
	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:$(DEST) product.proto