package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/lixd/grpc-go-example/features/proto/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

func unaryCall(c pb.EchoClient, requestID int, message string, want codes.Code) {
	// 每次都指定1秒超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.EchoRequest{Message: message}

	_, err := c.UnaryEcho(ctx, req)
	got := status.Code(err)
	fmt.Printf("[%v] wanted = %v, got = %v\n", requestID, want, got)
}

func streamingCall(c pb.EchoClient, requestID int, message string, want codes.Code) {
	// 每次都指定1秒超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := c.BidirectionalStreamingEcho(ctx)
	if err != nil {
		log.Printf("Stream err: %v", err)
		return
	}

	err = stream.Send(&pb.EchoRequest{Message: message})
	if err != nil {
		log.Printf("Send error: %v", err)
		return
	}

	_, err = stream.Recv()

	got := status.Code(err)
	fmt.Printf("[%v] wanted = %v, got = %v\n", requestID, want, got)
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewEchoClient(conn)

	unaryCall(c, 1, "world", codes.OK)
	unaryCall(c, 2, "delay", codes.DeadlineExceeded)
	unaryCall(c, 3, "[propagate me]world", codes.OK)
	unaryCall(c, 4, "[propagate me][propagate me]world", codes.DeadlineExceeded)
	streamingCall(c, 5, "[propagate me]world", codes.OK)
	streamingCall(c, 6, "[propagate me][propagate me]world", codes.DeadlineExceeded)
}
