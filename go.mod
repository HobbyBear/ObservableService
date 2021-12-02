module ObservableService

go 1.17

require google.golang.org/grpc v1.41.0

require (
	github.com/golang/protobuf v1.4.3
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4
)

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/gin-gonic/gin v1.7.7
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/client_golang v1.11.0
	github.com/uber/jaeger-client-go v2.29.1+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.uber.org/zap v1.19.1
	golang.org/x/text v0.3.5 // indirect
)
