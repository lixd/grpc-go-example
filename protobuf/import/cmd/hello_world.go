package main

import (
	"fmt"

	"github.com/lixd/grpc-go-example/protobuf/import"
)

func main() {
	c := proto.Computer{
		Name: "alienware",
		Cpu: &proto.CPU{
			Name:      "intel",
			Frequency: 4096,
		},
		Memory: &proto.Memory{
			Name: "芝奇",
			Cap:  8192,
		},
	}
	fmt.Println(c.String())
}
