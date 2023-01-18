package discovery

import (
	"context"

	"google.golang.org/grpc"

	discoveryv0 "github.com/aurae-runtime/ae/pkg/api/v0/discovery"
)

type Discovery interface {
	Discover(context.Context, *discoveryv0.DiscoverRequest) (*discoveryv0.DiscoverResponse, error)
}

type discovery struct {
	client discoveryv0.DiscoveryServiceClient
}

func New(ctx context.Context, conn grpc.ClientConnInterface) Discovery {
	return &discovery{
		client: discoveryv0.NewDiscoveryServiceClient(conn),
	}
}

func (d *discovery) Discover(ctx context.Context, req *discoveryv0.DiscoverRequest) (*discoveryv0.DiscoverResponse, error) {
	return d.client.Discover(ctx, req)
}
