package main

import (
	"context"
	"fmt"
	api "github.com/cacos-group/cacos/api"
	"google.golang.org/grpc"
	"time"
)

func main() {
	cc, err := grpc.Dial("127.0.0.1:8888", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := api.NewCacosClient(cc)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	reply, err := client.SayHello(ctx, &api.HelloRequest{
		Name: "adasdsa",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(reply, 1111)
}
