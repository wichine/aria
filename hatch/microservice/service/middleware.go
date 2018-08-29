package service

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

type Middleware func(ProductionService) ProductionService

func UseInstrumentingMiddleware(ints, chars metrics.Counter) Middleware {
	return func(next ProductionService) ProductionService {
		return instrumentingMiddleware{
			ints:  ints,
			chars: chars,
			next:  next,
		}
	}
}

type instrumentingMiddleware struct {
	ints  metrics.Counter
	chars metrics.Counter
	next  ProductionService
}

func (mw instrumentingMiddleware) AddProduction(ctx context.Context, prod *Production) (int64, error) {
	cnt, err := mw.next.AddProduction(ctx, prod)
	mw.ints.Add(float64(1))
	return cnt, err
}

func (mw instrumentingMiddleware) GetAllProduction(ctx context.Context) ([]Production, error) {
	pd, err := mw.next.GetAllProduction(ctx)
	mw.ints.Add(float64(1))
	return pd, err
}

func UseLoggingMiddleware(logger log.Logger) Middleware {
	return func(next ProductionService) ProductionService {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   ProductionService
}

func (mw loggingMiddleware) AddProduction(ctx context.Context, prod *Production) (int64, error) {
	defer func() {
		mw.logger.Log("method", "AddProduction")
	}()
	return mw.next.AddProduction(ctx, prod)
}

func (mw loggingMiddleware) GetAllProduction(ctx context.Context) ([]Production, error) {
	defer func() {
		mw.logger.Log("method", "GetAllProduction")
	}()
	return mw.next.GetAllProduction(ctx)
}
