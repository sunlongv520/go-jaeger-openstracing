package main

import (
	"google.golang.org/grpc"
	"go-opentracing/grpcpro/services"
	"net"
)

func main()  {
	rpcServer:=grpc.NewServer()
	services.RegisterProdServiceServer(rpcServer,new(services.ProdService))

	lis,_:=net.Listen("tcp",":8089")

	rpcServer.Serve(lis)


}
