package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/lixd/grpc-go-example/features/proto/echo"
	"google.golang.org/grpc"

	"google.golang.org/grpc/resolver"
)

const (
	myScheme      = "17x"
	myServiceName = "resolver.17x.lixueduan.com"

	backendAddr = "localhost:50051"
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

func makeRPCs(cc *grpc.ClientConn, n int) {
	hwc := pb.NewEchoClient(cc)
	for i := 0; i < n; i++ {
		callUnaryEcho(hwc, "this is examples/name_resolving")
	}
}

func main() {
	passthroughConn, err := grpc.Dial(
		// passthrough 也是gRPC内置的一个scheme
		fmt.Sprintf("passthrough:///%s", backendAddr), // Dial to "passthrough:///localhost:50051"
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer passthroughConn.Close()

	fmt.Printf("--- calling helloworld.Greeter/SayHello to \"passthrough:///%s\"\n", backendAddr)
	makeRPCs(passthroughConn, 10)

	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	exampleConn, err := grpc.DialContext(
		ctx,
		fmt.Sprintf("%s:///%s", myScheme, myServiceName), // Dial to "17x:///resolver.17x.lixueduan.com"
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer exampleConn.Close()

	fmt.Printf("--- calling helloworld.Greeter/SayHello to \"%s:///%s\"\n", myScheme, myServiceName)
	makeRPCs(exampleConn, 10)
}

// Following is an example name resolver. It includes a
// ResolverBuilder(https://godoc.org/google.golang.org/grpc/resolver#Builder)
// and a Resolver(https://godoc.org/google.golang.org/grpc/resolver#Resolver).
//
// A ResolverBuilder is registered for a scheme (in this example, "example" is
// the scheme). When a ClientConn is created for this scheme, the
// ResolverBuilder will be picked to build a Resolver. Note that a new Resolver
// is built for each ClientConn. The Resolver will watch the updates for the
// target, and send updates to the ClientConn.

// exampleResolverBuilder is a
// ResolverBuilder(https://godoc.org/google.golang.org/grpc/resolver#Builder).
type exampleResolverBuilder struct{}

func (*exampleResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &exampleResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			myServiceName: {backendAddr},
		},
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}
func (*exampleResolverBuilder) Scheme() string { return myScheme }

// exampleResolver is a
// Resolver(https://godoc.org/google.golang.org/grpc/resolver#Resolver).
type exampleResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *exampleResolver) ResolveNow(o resolver.ResolveNowOptions) {
	// 直接从map中取出对于的addrList
	addrStrs := r.addrsStore[r.target.Endpoint]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}

func (*exampleResolver) Close() {}

func init() {
	// Register the example ResolverBuilder. This is usually done in a package's
	// init() function.
	resolver.Register(&exampleResolverBuilder{})
}
