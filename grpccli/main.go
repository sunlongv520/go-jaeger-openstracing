package main

import (
	"context"
	"fmt"
	"go-opentracing/grpccli/until"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"go-opentracing/grpccli/services"
	"io"
	logger "log"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
)

func initJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		Sampler:&config.SamplerConfig{
			Type:     "const",
			Param:1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			//LocalAgentHostPort:  "192.168.1.234:6831",
			LocalAgentHostPort:  "192.168.1.234:6831",
		},
	}
	tracer, closer, err := cfg.New(service, config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("Error: connot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

func main(){

	tracer, closer := initJaeger("grpc-client")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	md := metadata.Pairs("key1","val1","key2","val2","key3","val3")
	ctx := metadata.NewOutgoingContext(context.Background(),md)


	span := tracer.StartSpan("say-hello")
	span.SetTag("hello-to", "helloTo")
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)



	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.Pairs()
	}

	if err := tracer.Inject(span.Context(), opentracing.HTTPHeaders,  until.MetadataTextMap(md)); err != nil {
		fmt.Println(ctx, "grpc_opentracing: failed serializing trace information: %v", err)
	}

	ctx = metadata.NewOutgoingContext(ctx, md)
	//ctx = metadata.AppendToOutgoingContext(ctx, util.TraceID, logs.GetTraceId(ctx))
	ctx = opentracing.ContextWithSpan(ctx, span)



	conn,err:=grpc.Dial(":8089",grpc.WithInsecure())
	if err!=nil{
		logger.Fatal(err)
	}
	defer conn.Close()

	prodClient:=services.NewProdServiceClient(conn)
	prodRes,err:=prodClient.GetProdStock(ctx,
		&services.ProdRequest{ProdId:12})
	if err!=nil{
		logger.Fatal(err)
	}
	fmt.Println(prodRes.ProdStock)
}
