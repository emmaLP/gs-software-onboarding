package grpc

import (
	"fmt"
	"net"

	pb "github.com/emmaLP/gs-software-onboarding/pkg/grpc/proto"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type server struct {
	port   int
	srv    pb.APIServer
	logger *zap.Logger
}

func NewServer(port int, logger *zap.Logger, srv pb.APIServer) *server {
	return &server{
		port:   port,
		srv:    srv,
		logger: logger,
	}
}

// Start the GRPC server
func (s *server) Start() error {
	s.logger.Debug(fmt.Sprintf("starting server on port %d", s.port))

	listenPort, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return errors.Wrap(err, "failed to listenPort")
	}

	gs := grpc.NewServer()
	pb.RegisterAPIServer(gs, s.srv)
	if err := gs.Serve(listenPort); err != nil {
		return errors.Wrap(err, "failed to serve")
	}

	return nil
}
