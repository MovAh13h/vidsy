PROTO_SRC = ./proto
OUTPUT = ./output
GO_OUT = go

.PHONY = all

all: build

build: proto $(OUTPUT)
	go build -o $(OUTPUT)/q.out ./cmd/q/main.go

proto: $(GO_OUT)
	protoc --proto_path=$(PROTO_SRC) --go_out=$(GO_OUT) --go-grpc_out=$(GO_OUT) $(PROTO_SRC)/*.proto

$(GO_OUT):
	@mkdir -p $@

$(OUTPUT):
	@mkdir -p $@

clean:
	rm -rf $(GO_OUT)
	rm -rf *.out
	rm -rf $(OUTPUT)