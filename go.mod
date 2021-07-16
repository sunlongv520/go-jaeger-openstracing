module go-opentracing

go 1.14

require (
	github.com/golang/protobuf v1.3.2
	github.com/micro/go-plugins/wrapper/trace/opentracing v0.0.0-20200119172437-4fe21aa238fd // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/uber/jaeger-client-go v2.29.1+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	google.golang.org/grpc v1.25.1
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
