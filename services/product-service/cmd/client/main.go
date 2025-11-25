package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/boretsotets/e-commerce-platform/product-service/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewProductServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.GetProduct(ctx, &pb.GetProductRequest{Id: 1})
	if err != nil {
		log.Fatalf("could not get product: %v", err)
	}
	fmt.Printf("Product: %#v\n", resp)
}
