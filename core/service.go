package core

import "github.com/go-kit/kit/endpoint"

type Service struct {
	Middleware []endpoint.Middleware
	Endpoint   endpoint.Endpoint
}

func (s *Service) AddMiddleware(ms ...endpoint.Middleware) {
	s.Middleware = append(ms, s.Middleware...)
}

func (s *Service) Compose() endpoint.Endpoint {
	final := s.Endpoint
	for _, m := range s.Middleware {
		final = m(final)
	}
	return final
}

func NewDefaultService() *Service {
	return &Service{
		Middleware: []endpoint.Middleware{},
	}
}
