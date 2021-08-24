package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/cacos-group/cacos/pkg/zaplog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"net/url"
	"sync"
	"time"
)

// ServerOption is an HTTP server option.
type ServerOption func(*Server)

// Address with server address.
func WithAddress(addr string) ServerOption {
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

// WithUnaryInterceptor returns a ServerOption that sets the UnaryServerInterceptor for the server.
func WithUnaryInterceptor(in ...grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.ints = in
	}
}

// WithOptions with grpc options.
func WithOptions(opts ...grpc.ServerOption) ServerOption {
	return func(s *Server) {
		s.grpcOpts = opts
	}
}

// Server is an GRPC server wrapper.
type Server struct {
	*grpc.Server
	ctx      context.Context
	lis      net.Listener
	tlsConf  *tls.Config
	once     sync.Once
	endpoint *url.URL
	err      error
	network  string
	address  string
	timeout  time.Duration
	log      zaplog.Logger

	ints     []grpc.UnaryServerInterceptor
	grpcOpts []grpc.ServerOption
}

// NewServer creates an HTTP server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network: "tcp",
		address: ":0",
		timeout: 1 * time.Second,
		log:     zaplog.DefaultLogger,
	}
	for _, o := range opts {
		o(srv)
	}
	var ints = []grpc.UnaryServerInterceptor{
		srv.unaryServerInterceptor(),
	}
	if len(srv.ints) > 0 {
		ints = append(ints, srv.ints...)
	}
	var grpcOpts = []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(ints...),
	}
	if srv.tlsConf != nil {
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(srv.tlsConf)))
	}
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}

	srv.Server = grpc.NewServer(grpcOpts...)
	return srv
}

// Endpoint return a real address to registry endpoint.
// examples:
//   grpc://127.0.0.1:9000?isSecure=false
func (s *Server) Endpoint() (*url.URL, error) {
	s.once.Do(func() {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			s.err = err
			return
		}
		//addr, err := host.Extract(s.address, lis)
		//if err != nil {
		//	_ = lis.Close()
		//	s.err = err
		//	return
		//}
		s.lis = lis
		//s.endpoint = endpoint.NewEndpoint("grpc", addr, s.tlsConf != nil)
	})
	if s.err != nil {
		return nil, s.err
	}
	return s.endpoint, nil
}

// Start start the gRPC server.
func (s *Server) Start(ctx context.Context) error {
	if _, err := s.Endpoint(); err != nil {
		return err
	}
	s.ctx = ctx
	s.log.Info(fmt.Sprintf("[gRPC] server listening on: %s", s.lis.Addr().String()))
	return s.Serve(s.lis)
}

// Stop stop the gRPC server.
func (s *Server) Stop(ctx context.Context) error {
	s.GracefulStop()
	s.log.Info("[gRPC] server stopping")
	return nil
}

func (s *Server) unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var (
			cancel context.CancelFunc
		)
		if s.timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, s.timeout)
			defer cancel()
		}

		resp, err := handler(ctx, req)
		return resp, err
	}
}
