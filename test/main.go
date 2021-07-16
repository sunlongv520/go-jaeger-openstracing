package main

import (
	"context"
	"io"
	"time"
	"fmt"


	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go/config"



)

/**
初始化
 */
func initJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		Sampler:&config.SamplerConfig{
			Type:     "const",
			Param:1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            false,
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

func TestDemo(req string, ctx context.Context) (reply string) {
	// 1. 创建span
	span, _ := opentracing.StartSpanFromContext(ctx, "span_testdemo_1")
	defer func() {
		// 4. 接口调用完，在tag中设置request和reply
		span.SetTag("request", req)
		span.SetTag("reply", reply)
		span.LogFields(
			log.String("event", "你又是谁？"),
			log.String("value", "我是你爷爷！^_^"),
		)
		span.Finish()
	}()

	println(req)
	//2. 模拟耗时
	time.Sleep(time.Second/2)
	//3. 返回reply
	reply = "TestDemoReply"
	return
}

// TestDemo2, 和上面TestDemo 逻辑代码一样
func TestDemo2(req string, ctx context.Context) (reply string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "span_testdemo2_1")
	defer func() {
		span.SetTag("request", req)
		span.SetTag("reply", reply)
		span.LogFields(
			log.String("event", "你是谁？"),
			log.String("value", "我是你爸！^_^"),
		)
		span.Finish()
	}()

	println(req)
	time.Sleep(time.Second/2)
	reply = "TestDemo2Reply"
	return
}

func main() {
	tracer, closer := initJaeger("jager-test-demo")
	defer closer.Close()
	//设置全局的tracer
	opentracing.SetGlobalTracer(tracer)
	//设置父的span
	span := tracer.StartSpan("span_root")


	ctx := opentracing.ContextWithSpan(context.Background(), span)

	r1 := TestDemo("Hello TestDemo", ctx)
	r2 := TestDemo2("Hello TestDemo2", ctx)

	fmt.Println(r1, r2)
	span.Finish()
}