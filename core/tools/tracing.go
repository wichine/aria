package tools

import (
	"aria/core/log"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	zipkinHttp "github.com/openzipkin/zipkin-go/reporter/http"
	"google.golang.org/grpc/metadata"
	"net/http"
	"strings"
)

var logger = log.GetLogger("AiraTools")
var globalTracer *zipkin.Tracer

func GetGlobalTracer() *zipkin.Tracer {
	if globalTracer != nil {
		return globalTracer
	}
	logger.Warningf("tracer not initialized, will use a noop one")
	noopTracer, _ := zipkin.NewTracer(nil, zipkin.WithNoopTracer(true))
	return noopTracer
}

func InitializeZipkin(reporterUrl, serviceName, localAddress string, on bool) error {
	opts := []zipkin.TracerOption{}
	reporter := zipkinHttp.NewReporter(reporterUrl)
	lep, err := zipkin.NewEndpoint(serviceName, localAddress)
	if err != nil {
		return fmt.Errorf("create zipkin endpoint error: %s", err)
	}
	opts = append(opts, zipkin.WithLocalEndpoint(lep))
	if !on {
		reporter = nil
		opts = append(opts, zipkin.WithNoopTracer(true), zipkin.WithNoopSpan(true))
	}
	globalTracer, err = zipkin.NewTracer(reporter, opts...)
	if err != nil {
		return fmt.Errorf("create zipkin tracer error: %s", err)
	}
	if on && !isZipkinApiValid(reporterUrl) {
		return fmt.Errorf("zipkin api %q not valid", reporterUrl)
	}
	return nil
}

func isZipkinApiValid(apiUrl string) bool {
	servicesUrl := strings.Replace(apiUrl, "spans", "services", 1)
	_, err := http.Get(servicesUrl)
	if err != nil {
		logger.Debugf("http get %q error: %s", servicesUrl, err)
		return false
	}
	return true
}

const spanKey = "zipkin-span-key-for-gincontext"

func GetZipkinSpanFromGinContext(ctx *gin.Context) zipkin.Span {
	if si, ok := ctx.Get(spanKey); ok {
		if span, ok := si.(zipkin.Span); ok {
			return span
		}
	}
	return GetGlobalTracer().StartSpan("")
}

func SetZipkinSpanToGinContext(ctx *gin.Context, span zipkin.Span) {
	ctx.Set(spanKey, span)
}

func SetZipkinSpanToGrpcHeader(ctx context.Context, span zipkin.Span) context.Context {
	md := &metadata.MD{}
	err := b3.InjectGRPC(md)(span.Context())
	if err != nil && err != b3.ErrEmptyContext {
		logger.Errorf("inject span to grpc metadata error: %s", err)
	}
	return metadata.NewOutgoingContext(ctx, *md)
}

// Get span from grpc context,create a new one if not found
func GetZipkinSpanFromGrpcHeader(ctx context.Context) zipkin.Span {
	spanOptions := []zipkin.SpanOption{}
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		parentSpan, err := b3.ExtractGRPC(&md)()
		if err != nil {
			logger.Debugf("extract span from context metadata error: %s", err)
		} else {
			spanOptions = append(spanOptions, zipkin.Parent(*parentSpan))
		}
	} else {
		logger.Debugf("no metadata found from incoming context.")
	}
	return GetGlobalTracer().StartSpan("", spanOptions...)
}

type grpcSpanKey struct{}

func GetZipkinSpanFromContext(ctx context.Context) zipkin.Span {
	if s, ok := ctx.Value(grpcSpanKey{}).(zipkin.Span); ok {
		return s
	}
	return GetGlobalTracer().StartSpan("")
}

func SetZipkinSpanToContext(ctx context.Context, span zipkin.Span) context.Context {
	return context.WithValue(ctx, grpcSpanKey{}, span)
}
