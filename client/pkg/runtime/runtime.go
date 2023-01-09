/* -------------------------------------------------------------------------- *\
 *             Apache 2.0 License Copyright © 2022 The Aurae Authors          *
 *                                                                            *
 *                +--------------------------------------------+              *
 *                |   █████╗ ██╗   ██╗██████╗  █████╗ ███████╗ |              *
 *                |  ██╔══██╗██║   ██║██╔══██╗██╔══██╗██╔════╝ |              *
 *                |  ███████║██║   ██║██████╔╝███████║█████╗   |              *
 *                |  ██╔══██║██║   ██║██╔══██╗██╔══██║██╔══╝   |              *
 *                |  ██║  ██║╚██████╔╝██║  ██║██║  ██║███████╗ |              *
 *                |  ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝ |              *
 *                +--------------------------------------------+              *
 *                                                                            *
 *                         Distributed Systems Runtime                        *
 *                                                                            *
 * -------------------------------------------------------------------------- *
 *                                                                            *
 *   Licensed under the Apache License, Version 2.0 (the "License");          *
 *   you may not use this file except in compliance with the License.         *
 *   You may obtain a copy of the License at                                  *
 *                                                                            *
 *       http://www.apache.org/licenses/LICENSE-2.0                           *
 *                                                                            *
 *   Unless required by applicable law or agreed to in writing, software      *
 *   distributed under the License is distributed on an "AS IS" BASIS,        *
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. *
 *   See the License for the specific language governing permissions and      *
 *   limitations under the License.                                           *
 *                                                                            *
\* -------------------------------------------------------------------------- */

package runtime

import (
	"context"
	"log"

	"google.golang.org/grpc"

	runtimev0 "github.com/aurae-runtime/ae/client/pkg/api/v0/runtime"
)

type Runtime interface {
	AllocateCell(context.Context, *runtimev0.CellServiceAllocateRequest) (*runtimev0.CellServiceAllocateResponse, error)
	FreeCell(context.Context, *runtimev0.CellServiceFreeRequest) error
	StartCell(context.Context, *runtimev0.CellServiceStartRequest) (*runtimev0.CellServiceStartResponse, error)
	StopCell(context.Context, *runtimev0.CellServiceStopRequest) error
}

type runtime struct {
	cellClient runtimev0.CellServiceClient
}

func New(ctx context.Context, conn grpc.ClientConnInterface) (Runtime, error) {
	r := &runtime{
		cellClient: runtimev0.NewCellServiceClient(conn),
	}

	return r, nil
}

func (r *runtime) AllocateCell(ctx context.Context, req *runtimev0.CellServiceAllocateRequest) (*runtimev0.CellServiceAllocateResponse, error) {
	rsp, err := r.cellClient.Allocate(ctx, req)
	if err != nil {
		log.Fatal("Cannot call Allocate ", err)

		return nil, err
	}

	return rsp, nil
}

func (r *runtime) FreeCell(ctx context.Context, req *runtimev0.CellServiceFreeRequest) error {
	_, err := r.cellClient.Free(ctx, req)
	if err != nil {
		log.Fatal("Cannot call Free ", err)

		return err
	}

	return nil
}

func (r *runtime) StartCell(ctx context.Context, req *runtimev0.CellServiceStartRequest) (*runtimev0.CellServiceStartResponse, error) {
	rsp, err := r.cellClient.Start(ctx, req)
	if err != nil {
		log.Fatal("Cannot call Start ", err)

		return nil, err
	}

	return rsp, nil
}

func (r *runtime) StopCell(ctx context.Context, req *runtimev0.CellServiceStopRequest) error {
	_, err := r.cellClient.Stop(ctx, req)
	if err != nil {
		log.Fatal("Cannot call Stop ", err)

		return err
	}

	return nil
}
