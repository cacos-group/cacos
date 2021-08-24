package main

import (
	"context"
	"fmt"
	api "github.com/cacos-group/cacos/api/gen/go"
	"google.golang.org/grpc"
	"time"
)

func main() {
	cc, err := grpc.Dial("127.0.0.1:9090", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := api.NewCacosClient(cc)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	//reply, err := client.AuthLogin(ctx, &api.LoginRequest{
	//	Username: "adasdsa",
	//	Password: "ssd",
	//})
	//reply, err := client.NamespaceList(ctx, &api.NamespaceListReq{})
	//if err != nil {
	//	panic(err)
	//}

	reply, err := client.AppList(ctx, &api.AppListReq{
		Namespace: "namespace5",
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", reply.AppList)
}
