package service

import (
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/pkg/errors"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

type Service interface {
	Run() error
	Stop() error
}

type CommonService struct {
	app *kratos.App
}

func NewService(logger log.Logger, hs *http.Server, gs *grpc.Server, options ...kratos.Option) Service {
	options = append(options, kratos.Metadata(map[string]string{}),
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Logger(log.DefaultLogger),
		kratos.Server(
			hs,
			gs,
		))
	app := kratos.New(
		options...,
	)

	return &CommonService{
		app: app,
	}
}

func (s CommonService) Run() error {
	// start and wait for stop signal
	return s.app.Run()
}

func (s CommonService) Stop() error {
	return errors.WithMessage(nil, "Unimplemented")
}
