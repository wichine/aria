package core

import (
	"aria/hatch/microservice/core/config"
	"aria/hatch/microservice/core/log"
	"aria/hatch/microservice/core/svcdiscovery"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"google.golang.org/grpc"
	"net"
	"strings"
)

type Aria struct {
	Config       *config.AriaConfig
	GrpcServer   *grpc.Server
	GrpcListener net.Listener
	HttpEngine   *gin.Engine
	logger       *logging.Logger
}

func New(config *config.AriaConfig) *Aria {
	port := strings.Split(config.Address, ":")[1]
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	logger := &logging.Logger{
		Module: "@aria_core",
	}
	loggerBackend := logging.NewBackendFormatter(
		logging.NewLogBackend(log.DefaultLogWriter, "", 0),
		logging.MustStringFormatter(`%{color}[ARIA][%{time:2006-01-02 15:04:05.000}] %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`),
	)
	levelBackEnd := logging.AddModuleLevel(loggerBackend)
	logger.SetBackend(levelBackEnd)
	a := &Aria{
		Config:       config,
		GrpcListener: lis,
		GrpcServer:   grpc.NewServer(),
		logger:       logger,
		//log: func(format string, i ...interface{}) {
		//	pre := fmt.Sprintf("[ARIA] %v | ", time.Now().Format("2006/01/02 - 15:04:05.000"))
		//	fmt.Fprintf(log.DefaultLogWriter, pre+format, i)
		//},
	}
	a.printConfig()
	return a
}

func (a *Aria) RegisterAll(ts ...Transport) {
	for _, s := range ts {
		s.Register(a.GrpcServer)
	}
}
func (a *Aria) ServeGRPC() error {
	if err := a.GrpcServer.Serve(a.GrpcListener); err != nil {
		return err
	}
	return nil
}

func (a *Aria) Run() error {
	// register service
	if a.Config.ServiceDiscovery.Enable {
		a.logger.Debugf("Start %s service discovery ", "etcd")
		// register to etcd service discovery center
		s, err := svcdiscovery.GetEtcdServiceDiscoveryInstance()
		if err != nil {
			panic(err)
		}
		a.logger.Debugf("Register service [%s] at [%s]", a.Config.ServiceKey, a.Config.Address)
		if err = s.Register(a.Config.ServiceKey, a.Config.Address); err != nil {
			panic(err)
		}
		defer s.DeRegister()
	}
	// serve grpc / http
	a.logger.Debugf("Start service at %s", a.Config.Address)
	if err := a.GrpcServer.Serve(a.GrpcListener); err != nil {
		return err
	}
	return nil
}

func (a *Aria) printConfig() {
	params := Flatten(a.Config)
	var buffer bytes.Buffer
	for i := range params {
		buffer.WriteString("\n\t")
		buffer.WriteString(params[i])
	}
	a.logger.Debugf("Aria config: %s\n", buffer.String())
}

type Transport interface {
	Register(s *grpc.Server)
}
