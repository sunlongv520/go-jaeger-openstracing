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

func TestDemo(req string, ctx context.Context) (reply string) {
	// 1. 创建span
	span, _ := opentracing.StartSpanFromContext(ctx, "span_testdemo_sl")
	defer span.Finish()
	span.SetTag("request", req)
	span.SetTag("reply", "TestDemo")
	span.LogFields(
		log.String("event", "你又是谁？"),
		log.String("value", "TestDemo^_^"),
	)
	fmt.Printf("进入函数:TestDemo，spanid：%s \r\n",span)
	//2. 模拟耗时
	time.Sleep(time.Second/2)
	//3. 返回reply
	reply = "回复：TestDemoReply \r\n"
	return
}

// TestDemo2, 和上面TestDemo 逻辑代码一样
func TestDemo2(req string, ctx context.Context) (reply string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "span_testdemo2_sl")
	defer span.Finish()

	span.SetTag("request", req)
	span.SetTag("reply", "TestDemo2")
	span.LogFields(
		log.String("event", "你是谁？"),
		log.String("value", "我是TestDemo2^_^"),
	)

	fmt.Printf("进入函数:TestDemo2,spanid:%s \r\n",span)
	time.Sleep(time.Second/2)
	reply = "回复：TestDemo2Reply \r\n"
	ctx2 := opentracing.ContextWithSpan(ctx, span)
	_ = TestDemo3("Hello TestDemo3", ctx2)
	//fmt.Println(r3)
	return
}

func TestDemo3(req string, ctx context.Context) (reply string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "span_testdemo3_sl")
	defer span.Finish()

	span.SetTag("request", req)
	span.SetTag("reply", "TestDemo3")
	span.LogFields(
		log.String("event", "你是谁？"),
		log.String("value", "TestDemo2的儿子^_^"),
	)

	span.LogKV("event2", "println")
	fmt.Printf("进入函数:TestDemo3,spanid:%s \r\n",span)
	time.Sleep(time.Second/2)
	reply = "TestDemo2Reply"

	ctx2 := opentracing.ContextWithSpan(ctx, span)
	_ = TestDemo4("Hello TestDemo3", ctx2)
	return
}

func TestDemo4(req string, ctx context.Context) (reply string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "span_testdemo4_sl")
	defer span.Finish()

	span.SetTag("request", req)
	span.SetTag("reply", "TestDemo4" +
		"+")
	span.LogFields(
		log.String("event", "你是谁？"),
		log.String("value", "TestDemo3的儿子^_^"),
	)

	span.LogKV("event2", "println")

	fmt.Printf("进入函数:TestDemo4,spanid:%s \r\n",span)
	time.Sleep(time.Second/2)
	reply = "回复：TestDemo2Reply \r\n"
	return
}

func main() {
	tracer, closer := initJaeger("jager-test-function")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	span := tracer.StartSpan("span_root_sl")
	fmt.Printf("父SpanId:%s \r\n",span)

	ctx := opentracing.ContextWithSpan(context.Background(), span)
	r1 := TestDemo("Hello TestDemo", ctx)
	r2 := TestDemo2("Hello TestDemo2", ctx)
	fmt.Println(r1, r2)
	span.Finish()

}