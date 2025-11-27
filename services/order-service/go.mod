module github.com/boretsotets/e-commerce-platform/services/order-service

go 1.24.0

toolchain go1.24.9

require (
	google.golang.org/grpc v1.76.0
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/confluentinc/confluent-kafka-go/v2 v2.12.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
)

require github.com/boretsotets/e-commerce-platform/product-service v0.0.0

replace github.com/boretsotets/e-commerce-platform/product-service => ../product-service
