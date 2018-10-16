package middleware

import (
	"aria/core"
	"aria/core/tools"
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/op/go-logging"
	"google.golang.org/grpc"
	"path"
	"reflect"
)

func WithMiddleware(serviceImpl interface{}, mws ...endpoint.Middleware) core.Transport {
	sPtr := reflect.ValueOf(serviceImpl)
	if sPtr.Kind().String() != "ptr" {
		panic("serviceImpl should be a ptr of a service implement object!")
	}
	ariaTransportType := reflect.TypeOf((*core.Transport)(nil)).Elem()
	if !sPtr.Type().Implements(ariaTransportType) {
		panic(fmt.Sprintf("serviceImpl type need to be %q but %q given", ariaTransportType.String(), sPtr.Type().String()))
	}
	s := sPtr.Elem()

	for i := 0; i < s.NumField(); i++ {
		if m := s.Field(i).MethodByName("AddMiddleware"); m.IsValid() {
			in := []reflect.Value{}
			for _, mw := range mws {
				in = append(in, reflect.ValueOf(mw))
			}
			m.Call(in)
		}
	}
	return serviceImpl.(core.Transport)
}

func LogMiddleware(logger *logging.Logger) endpoint.Middleware {
	return func(ep endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			response, err = ep(ctx, request)
			if err != nil {
				logger.Errorf("[Middleware Log] mothod: %q, request: %v, error: %q", grpc.ServerTransportStreamFromContext(ctx).Method(), request, err)
				return
			}
			logger.Debugf("[Middleware Log] mothod: %q, request: %v, response: %v", grpc.ServerTransportStreamFromContext(ctx).Method(), request, response)
			return
		}
	}
}

func ZipkinMiddleware(ep endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, _ error) {
		rawMethod := grpc.ServerTransportStreamFromContext(ctx).Method()
		method := path.Base(rawMethod)
		span := tools.GetZipkinSpanFromGrpcHeader(ctx)
		span.SetName(fmt.Sprintf("grpc server: %s", method))
		defer span.Finish()
		newCtx := tools.SetZipkinSpanToGrpcHeader(tools.SetZipkinSpanToContext(ctx, span), span)
		return ep(newCtx, request)
	}
}
