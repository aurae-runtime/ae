package observe

import (
	"context"

	"google.golang.org/grpc"

	observev0 "github.com/aurae-runtime/ae/pkg/api/v0/observe"
)

type Observe interface {
	GetAuraeDaemonLogStream(context.Context, *observev0.GetAuraeDaemonLogStreamRequest) (observev0.ObserveService_GetAuraeDaemonLogStreamClient, error)
	GetSubProcessStream(context.Context, *observev0.GetSubProcessStreamRequest) (observev0.ObserveService_GetSubProcessStreamClient, error)
}

type observe struct {
	client observev0.ObserveServiceClient
}

func New(ctx context.Context, conn grpc.ClientConnInterface) Observe {
	return &observe{
		client: observev0.NewObserveServiceClient(conn),
	}
}

func (o *observe) GetAuraeDaemonLogStream(ctx context.Context, req *observev0.GetAuraeDaemonLogStreamRequest) (observev0.ObserveService_GetAuraeDaemonLogStreamClient, error) {
	return o.client.GetAuraeDaemonLogStream(ctx, req)
}

func (o *observe) GetSubProcessStream(ctx context.Context, req *observev0.GetSubProcessStreamRequest) (observev0.ObserveService_GetSubProcessStreamClient, error) {
	return o.client.GetSubProcessStream(ctx, req)
}
