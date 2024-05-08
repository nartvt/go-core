package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/nartvt/go-core/conf"
	"github.com/nartvt/go-core/middleware/jwt"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, logger log.Logger) *khttp.Server {
	authMiddleware := jwt.Server()
	if c.Auth != nil {
		authMiddleware = jwt.Server(jwt.WithRequired(c.Auth.Required), jwt.WithExcludes(c.Auth.Excludes), jwt.WithAutoParse(c.Auth.AutoParse))
	}
	var opts = []khttp.ServerOption{
		khttp.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			logging.Server(logger),
			metrics.Server(),
			validate.Validator(),
			authMiddleware,
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, khttp.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, khttp.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, khttp.Timeout(c.Http.Timeout.AsDuration()))
	}

	filter := khttp.Filter(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if request.RequestURI == "/health" {
				writer.Header().Add("Content-Type", "application/json")
				_ = json.NewEncoder(writer).Encode(map[string]bool{"ok": true})
			} else {
				handler.ServeHTTP(writer, request)
			}
		})
	})
	opts = append(opts, filter)
	srv := khttp.NewServer(opts...)
	//v1.RegisterGreeterHTTPServer(srv, greeter)

	//_metricSeconds := prometheus.NewHistogramVec(prometheus.HistogramOpts{
	//	Namespace: "server",
	//	Subsystem: "requests",
	//	Name:      "duration_ms",
	//	Help:      "server requests duration(ms).",
	//	Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	//}, []string{"kind", "operation"})
	//
	//_metricRequests := prometheus.NewCounterVec(prometheus.CounterOpts{
	//	Namespace: "client",
	//	Subsystem: "requests",
	//	Name:      "code_total",
	//	Help:      "The total number of processed requests",
	//}, []string{"kind", "operation", "code", "reason"})
	//
	//prometheus.MustRegister(_metricSeconds, _metricRequests)
	//
	//srv.Handle("/metrics", promhttp.Handler())

	return srv
}
