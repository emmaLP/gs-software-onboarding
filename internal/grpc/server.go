package grpc

import (
	"fmt"
	"net"

	pb "github.com/emmaLP/gs-software-onboarding/pkg/grpc/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type server struct {
	port   int
	srv    pb.APIServer
	logger *zap.Logger
}

// NewServer instantiates the struct for use when starting the GRPC server
func NewServer(port int, logger *zap.Logger, srv pb.APIServer) *server {
	return &server{
		port:   port,
		srv:    srv,
		logger: logger,
	}
}

// Start the GRPC server
func (s *server) Start() (*grpc.Server, error) {
	s.logger.Debug(fmt.Sprintf("starting server on port %d", s.port))

	listenPort, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port, %w", err)
	}

	gs := grpc.NewServer()
	pb.RegisterAPIServer(gs, s.srv)
	if err := gs.Serve(listenPort); err != nil {
		return nil, fmt.Errorf("failed to serve. %w", err)
	}

	return gs, nil
}
