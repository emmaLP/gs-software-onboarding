package grpc

import (
	"context"
	"fmt"
	"io"

	"github.com/emmaLP/gs-software-onboarding/pkg/common/model"
	pb "github.com/emmaLP/gs-software-onboarding/pkg/grpc/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client interface {
	ListAll(ctx context.Context) ([]*model.Item, error)
	ListStories(ctx context.Context) ([]*model.Item, error)
	ListJobs(ctx context.Context) ([]*model.Item, error)
}

type client struct {
	grpcClient     pb.APIClient
	grpcConnection *grpc.ClientConn
	logger         *zap.Logger
}

// NewClient instantiates a connection to a grpc server
func NewClient(addr string, logger *zap.Logger) (*client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("connecting to grpc server with address %s. Error: %w", addr, err)
	}
	logger.Debug("Connected to GRPC server successfully")

	apiClient := pb.NewAPIClient(conn)
	logger.Debug("GRPC Client instantiated")

	return &client{
		grpcClient:     apiClient,
		grpcConnection: conn,
		logger:         logger,
	}, nil
}

func (c *client) ListAll(ctx context.Context) ([]*model.Item, error) {
	stream, err := c.grpcClient.ListAll(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("An error occurred when streaming all. %w", err)
	}
	return handleStreamItems(ctx, stream)
}

func (c *client) ListStories(ctx context.Context) ([]*model.Item, error) {
	stream, err := c.grpcClient.ListStories(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("An error occurred when streaming stories. %w", err)
	}
	return handleStreamItems(ctx, stream)
}

func (c *client) ListJobs(ctx context.Context) ([]*model.Item, error) {
	stream, err := c.grpcClient.ListJobs(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("An error occurred when streaming jobs. %w", err)
	}
	return handleStreamItems(ctx, stream)
}

func (c *client) Close() {
	err := c.grpcConnection.Close()
	if err != nil {
		c.logger.Error("Unable to close connection", zap.Error(err))
	}
}

func handleStreamItems(ctx context.Context, stream interface{ Recv() (*pb.Item, error) }) ([]*model.Item, error) {
	var items []*model.Item
	isComplete := false
	for !isComplete {
		select {
		case <-ctx.Done():
			isComplete = true
		default:
			pbItem, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					isComplete = true
					break
				} else {
					return nil, fmt.Errorf("receiving item from server. %w", err)
				}
			}
			item := model.PItemToItem(pbItem)
			items = append(items, &item)
		}
	}
	return items, nil
}
