package discovery

import (
	"context"

	"google.golang.org/grpc"

	discoveryv0 "github.com/aurae-runtime/ae/pkg/api/v0/discovery"
)

type Discovery interface {
	Health(context.Context, *discoveryv0.HealthRequest) (*discoveryv0.HealthResponse, error)
}

type discovery struct {
	client discoveryv0.DiscoveryServiceClient
}

func New(ctx context.Context, conn grpc.ClientConnInterface) Discovery {
	return &discovery{
		client: discoveryv0.NewDiscoveryServiceClient(conn),
	}
}

func (d *discovery) Health(ctx context.Context, req *discoveryv0.HealthRequest) (*discoveryv0.HealthResponse, error) {
	return d.client.Health(ctx, req)
}
