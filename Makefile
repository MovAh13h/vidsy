SRC = ./proto
OUT = output
GO_OUT = go

.PHONY = all

all: | $(OUT) $(GO_OUT)
	protoc --proto_path=$(SRC) --cpp_out=$(OUT)/cpp --go_out=$(GO_OUT) --go-grpc_out=$(GO_OUT) $(SRC)/*.proto

$(OUT):
	mkdir -p $@
	mkdir -p $@/cpp

$(GO_OUT):
	mkdir -p $@

clean:
	rm -rf $(OUT)
	rm -rf $(GO_OUT)
