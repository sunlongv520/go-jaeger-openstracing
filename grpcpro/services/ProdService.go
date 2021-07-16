package services

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"io"
	"go-opentracing/grpcpro/until"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
)

type ProdService struct {

}

func(this *ProdService) GetProdStock(ctx context.Context, request *ProdRequest) (*ProdResponse, error) {
		md,ok :=metadata.FromIncomingContext(ctx)
		fmt.Println(md,ok)
		if !ok {
			md = metadata.New(nil)
		}
		tracer, closer := initJaeger("grpc-server")
		defer closer.Close()
		spanContext, err := tracer.Extract(opentracing.HTTPHeaders,until.MetadataTextMap(md))
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			grpclog.Errorf("extract from metadata err %v", err)
		}
		//开始追踪该方法
		serverSpan := tracer.StartSpan(
			"grpc-server-name",
			ext.RPCServerOption(spanContext),
			ext.SpanKindRPCServer,
		)
		serverSpan.SetTag("grpc-server-tag", "grpc-server-tag-value")
		ctx = opentracing.ContextWithSpan(ctx, serverSpan)

		defer serverSpan.Finish()

	  	return &ProdResponse{ProdStock:20},nil
}


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
