package service

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"net/http"
	"os"
	"prometheus/client_golang/prometheus/promhttp"
)

type production struct {
	production []Production
}

func NewProductionService() ProductionService {
	p := NewBasicProductionService()
	ints := prometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "example",
		Subsystem: "addsvc",
		Name:      "integers_summed",
		Help:      "Total count of integers summed via the Sum method.",
	}, []string{})
	chars := prometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "example",
		Subsystem: "addsvc",
		Name:      "characters_concatenated",
		Help:      "Total count of characters concatenated via the Concat method.",
	}, []string{})
	p = UseInstrumentingMiddleware(ints, chars)(p)
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	p = UseLoggingMiddleware(logger)(p)
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":8080", nil)
	return p
}

func NewBasicProductionService() ProductionService {
	return &production{[]Production{}}
}
func (p *production) AddProduction(ctx context.Context, prod *Production) (int64, error) {
	p.production = append(p.production, *prod)
	return int64(len(p.production)), nil
}

func (p *production) GetAllProduction(ctx context.Context) ([]Production, error) {
	return p.production, nil
}
