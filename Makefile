.PHONY: proto
proto:
	protoc --proto_path=./model \
		--go_out=paths=source_relative:./model \
		./model/unollm.proto
