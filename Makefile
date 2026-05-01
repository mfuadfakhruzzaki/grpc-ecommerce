PROTO_DIR=proto
OUT_DIR=proto

.PHONY: proto
proto:
	protoc \
		-I$(PROTO_DIR) \
		-I$(PROTO_DIR)/google \
		--go_out=$(OUT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) \
		--go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(OUT_DIR) \
		--grpc-gateway_opt=paths=source_relative \
		$(PROTO_DIR)/user/v1/user.proto \
		$(PROTO_DIR)/product/v1/product.proto \
		$(PROTO_DIR)/order/v1/order.proto

.PHONY: up
up:
	docker compose up --build

.PHONY: down
down:
	docker compose down -v