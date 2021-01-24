package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/lixd/grpc-go-example/data"
	ecpb "github.com/lixd/grpc-go-example/features/proto/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

func main() {
	// 加载客户端证书
	certificate, err := tls.LoadX509KeyPair(data.Path("x509/client.crt"), data.Path("x509/client.key"))
	if err != nil {
		log.Fatal(err)
	}
	// 构建CertPool以校验服务端证书有效性
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(data.Path("x509/ca.crt"))
	if err != nil {
		log.Fatal(err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatal("failed to append ca certs")
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certificate},
		ServerName:   "www.lixueduan.com", // NOTE: this is required!
		RootCAs:      certPool,
	})
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("DialContext error:%v", err)
	}
	defer conn.Close()
	client := ecpb.NewEchoClient(conn)
	callUnaryEcho(client, "hello world")
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
