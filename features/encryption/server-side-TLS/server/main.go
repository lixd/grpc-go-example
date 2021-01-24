/*
 *
 * Copyright 2018 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Binary server is an example server.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/lixd/grpc-go-example/data"
	ecpb "github.com/lixd/grpc-go-example/features/proto/echo"
)

var port = flag.Int("port", 50051, "the port to serve on")

type ecServer struct {
	ecpb.UnimplementedEchoServer
}

func (s *ecServer) UnaryEcho(ctx context.Context, req *ecpb.EchoRequest) (*ecpb.EchoResponse, error) {
	return &ecpb.EchoResponse{Message: req.Message}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 指定使用服务端证书创建一个 TLS credentials。
	creds, err := credentials.NewServerTLSFromFile(data.Path("x509/server.crt"), data.Path("x509/server.key"))
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}
	// 指定使用 TLS credentials。
	s := grpc.NewServer(grpc.Creds(creds))
	ecpb.RegisterEchoServer(s, &ecServer{})
	log.Println("Server gRPC on 0.0.0.0" + fmt.Sprintf(":%d", *port))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
