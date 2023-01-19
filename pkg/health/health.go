package health

import (
	"context"

	"google.golang.org/grpc"

	//healthv1 "github.com/aurae-runtime/ae/pkg/api/grpc/health/v1/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

type Health interface {
	Check(context.Context, *healthv1.HealthCheckRequest) (*healthv1.HealthCheckResponse, error)
}

type health struct {
	client healthv1.HealthClient
}

func New(ctx context.Context, conn grpc.ClientConnInterface) Health {
	return &health{
		client: healthv1.NewHealthClient(conn),
	}
}

func (h *health) Check(ctx context.Context, req *healthv1.HealthCheckRequest) (*healthv1.HealthCheckResponse, error) {
	return h.client.Check(ctx, req)
}
