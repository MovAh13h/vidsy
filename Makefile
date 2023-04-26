PROTO_SRC = ./proto
OUTPUT = ./output
GO_OUT = go

.PHONY = all

all: build

build: proto $(OUTPUT)
	go build -o $(OUTPUT)/consumer.out ./cmd/consumer/main.go ./cmd/consumer/consumer.go
	go build -o $(OUTPUT)/q.out ./cmd/q/main.go ./cmd/q/video_server.go
	go build -o $(OUTPUT)/client.out ./cmd/client/main.go

proto: $(GO_OUT)
	protoc --proto_path=$(PROTO_SRC) --go_out=$(GO_OUT) --go-grpc_out=$(GO_OUT) $(PROTO_SRC)/*.proto

run_q:
	$(OUTPUT)/q.out

run_consumer:
	$(OUTPUT)/consumer.out

run_client:
	$(OUTPUT)/client.out

$(GO_OUT):
	@mkdir -p $@

$(OUTPUT):
	@mkdir -p $@

clean:
	rm -rf $(GO_OUT)
	rm -rf *.out
	rm -rf $(OUTPUT)