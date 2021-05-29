package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/lixd/grpc-go-example/features/proto/echo"
	"github.com/sercand/kuberesolver/v3"
	"google.golang.org/grpc"
)

func callUnaryEcho(c pb.EchoClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.UnaryEcho(ctx, &pb.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Println(r.Message)
}

func main() {
	// Register kuberesolver to grpc
	kuberesolver.RegisterInCluster()

	// Make ClientConn with round_robin policy.
	cc, err := grpc.Dial(
		// fmt.Sprintf("kubernetes:///service.namespace:portname"),
		fmt.Sprintf("kubernetes:///svc-mygrpc.myns:50051"),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), // This sets the initial balancing policy.
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer cc.Close()

	hwc := pb.NewEchoClient(cc)
	for i := 0; i < 10; i++ {
		callUnaryEcho(hwc, "this is examples/load_balancing")
	}
}
