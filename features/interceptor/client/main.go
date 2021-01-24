package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	ecpb "github.com/lixd/grpc-go-example/features/proto/echo"
	"golang.org/x/oauth2"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/oauth"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

const fallbackToken = "some-secret-token"

// logger 简单打印日志
func logger(format string, a ...interface{}) {
	fmt.Printf("LOG:\t"+format+"\n", a...)
}

// unaryInterceptor 一个简单的 unary interceptor 示例。
func unaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	var credsConfigured bool
	for _, o := range opts {
		_, ok := o.(grpc.PerRPCCredsCallOption)
		if ok {
			credsConfigured = true
			break
		}
	}
	if !credsConfigured {
		opts = append(opts, grpc.PerRPCCredentials(oauth.NewOauthAccess(&oauth2.Token{
			AccessToken: fallbackToken,
		})))
	}
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	end := time.Now()
	logger("RPC: %s, start time: %s, end time: %s, err: %v", method, start.Format("Basic"), end.Format(time.RFC3339), err)
	return err
}

// wrappedStream  用于包装 grpc.ClientStream 结构体并拦截其对应的方法。
type wrappedStream struct {
	grpc.ClientStream
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	logger("Receive a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	logger("Send a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.SendMsg(m)
}

func newWrappedStream(s grpc.ClientStream) grpc.ClientStream {
	return &wrappedStream{s}
}

// streamInterceptor 一个简单的 stream interceptor 示例。
func streamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	var credsConfigured bool
	for _, o := range opts {
		_, ok := o.(*grpc.PerRPCCredsCallOption)
		if ok {
			credsConfigured = true
			break
		}
	}
	if !credsConfigured {
		opts = append(opts, grpc.PerRPCCredentials(oauth.NewOauthAccess(&oauth2.Token{
			AccessToken: fallbackToken,
		})))
	}
	s, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		return nil, err
	}
	return newWrappedStream(s), nil
}

func callUnaryEcho(client ecpb.EchoClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.UnaryEcho(ctx, &ecpb.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("client.UnaryEcho(_) = _, %v: ", err)
	}
	fmt.Println("UnaryEcho: ", resp.Message)
}

func callBidiStreamingEcho(client ecpb.EchoClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c, err := client.BidirectionalStreamingEcho(ctx)
	if err != nil {
		return
	}
	for i := 0; i < 5; i++ {
		if err := c.Send(&ecpb.EchoRequest{Message: fmt.Sprintf("Request %d", i+1)}); err != nil {
			log.Fatalf("failed to send request due to error: %v", err)
		}
	}
	err = c.CloseSend()
	if err != nil {
		log.Printf("close send err:%v", err)
	}
	for {
		resp, err := c.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {

			log.Fatalf("failed to receive response due to error: %v", err)
		}
		fmt.Println("BidiStreaming Echo: ", resp.Message)
	}
}

func main() {
	flag.Parse()

	// Create tls based credential.
	// creds, err := credentials.NewClientTLSFromFile(data.Path("x509/server.crt"), "www.lixueduan.com")
	// if err != nil {
	// 	log.Fatalf("failed to load credentials: %v", err)
	// }

	// Set up a connection to the server.
	// conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(creds), grpc.WithUnaryInterceptor(unaryInterceptor), grpc.WithStreamInterceptor(streamInterceptor), grpc.WithBlock())
	conn, err := grpc.Dial(*addr /*grpc.WithTransportCredentials(creds),*/, grpc.WithInsecure(), grpc.WithUnaryInterceptor(unaryInterceptor), grpc.WithStreamInterceptor(streamInterceptor), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Make a echo client and send RPCs.
	client := ecpb.NewEchoClient(conn)
	callUnaryEcho(client, "hello world")
	callBidiStreamingEcho(client)
}
