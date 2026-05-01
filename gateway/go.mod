module github.com/mfuadfakhruzzaki/grpc-ecommerce/gateway

go 1.25.5

replace github.com/mfuadfakhruzzaki/grpc-ecommerce => ../

require (
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.29.0
	github.com/mfuadfakhruzzaki/grpc-ecommerce v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.80.0
)

require (
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260414002931-afd174a4e478 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260414002931-afd174a4e478 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
