package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/lixd/grpc-go-example/features/proto/echo"
	"google.golang.org/grpc"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

func sendMessage(stream pb.Echo_BidirectionalStreamingEchoClient, msg string) error {
	fmt.Printf("sending message %q\n", msg)
	return stream.Send(&pb.EchoRequest{Message: msg})
}

func recvMessage(stream pb.Echo_BidirectionalStreamingEchoClient) {
	res, err := stream.Recv()
	if err != nil {
		fmt.Printf("stream.Recv() returned error %v\n", err)
		return
	}
	fmt.Printf("received message %q\n", res.GetMessage())
}

func main() {
	flag.Parse()

	// 建立连接
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewEchoClient(conn)

	// 初始化一个带取消功能的ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	stream, err := c.BidirectionalStreamingEcho(ctx)
	if err != nil {
		log.Fatalf("error creating stream: %v", err)
	}

	// 正常发送消息
	if err := sendMessage(stream, "hello"); err != nil {
		log.Fatalf("error sending on stream: %v", err)
	}
	if err := sendMessage(stream, "world"); err != nil {
		log.Fatalf("error sending on stream: %v", err)
	}

	// 正常接收消息
	recvMessage(stream)
	recvMessage(stream)
	// 这里调用cancel方法取消 ctx
	fmt.Println("cancelling context")
	cancel()
	time.Sleep(time.Second)
	// 再次发送消息 这里是否会报错取决于ctx是否检测到前面发送的取消命令(cancel())
	if err := sendMessage(stream, "world"); err != nil {
		log.Printf("error sending on stream: %v", err)
	}

	// 这里一定会报错
	recvMessage(stream)
}
