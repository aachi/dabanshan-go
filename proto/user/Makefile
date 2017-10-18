DEST = ./
PROTOBUFS = $(wildcard *.proto)
PROTOC_FLAGS  = -I .

all:
	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:$(DEST) ProductModel.proto
	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:$(DEST) ProductService.proto
	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:$(DEST) UserModel.proto
	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:$(DEST) UserService.proto