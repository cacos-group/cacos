package http

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/cacos-group/cacos/pkg/zaplog"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// ServerOption is an HTTP server option.
type ServerOption func(*Server)

// Address with server address.
func Address(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

// WithLogger with server logger.
func WithLogger(logger zaplog.Logger) ServerOption {
	return func(s *Server) {
		s.log = logger
	}
}

// Server is an HTTP server wrapper.
type Server struct {
	*http.Server
	lis      net.Listener
	tlsConf  *tls.Config
	once     sync.Once
	endpoint *url.URL
	err      error
	network  string
	address  string
	timeout  time.Duration
	//filters  []FilterFunc
	//ms       []middleware.Middleware
	//dec      DecodeRequestFunc
	//enc      EncodeResponseFunc
	//ene      EncodeErrorFunc
	//router   *mux.Router
	log zaplog.Logger
}

// NewServer creates an HTTP server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network: "tcp",
		address: ":0",
		timeout: 1 * time.Second,
		//dec:     DefaultRequestDecoder,
		//enc:     DefaultResponseEncoder,
		//ene:     DefaultErrorEncoder,
		log: zaplog.DefaultLogger,
	}
	for _, o := range opts {
		o(srv)
	}
	srv.Server = &http.Server{
		//Handler:   FilterChain(srv.filters...)(srv),
		TLSConfig: srv.tlsConf,
	}
	//srv.router = mux.NewRouter()
	//srv.router.Use(srv.filter())
	return srv
}

// Start start the HTTP server.
func (s *Server) Start(ctx context.Context) error {
	//if _, err := s.Endpoint(); err != nil {
	//	return err
	//}
	s.BaseContext = func(net.Listener) context.Context {
		return ctx
	}
	s.log.Info(fmt.Sprintf("[HTTP] server listening on: %s", s.lis.Addr().String()))
	var err error
	if s.tlsConf != nil {
		err = s.ServeTLS(s.lis, "", "")
	} else {
		err = s.Serve(s.lis)
	}
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Stop stop the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("[HTTP] server stopping")
	return s.Shutdown(ctx)
}
