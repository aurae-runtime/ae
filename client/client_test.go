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

package client

import (
	"context"
	"testing"

	runtimev0 "github.com/aurae-runtime/ae/client/pkg/api/v0/runtime"
)

const (
	testCellName = "MyBrandNewCell"
)

func TestClientRuntimeAllocate(t *testing.T) {
	c, err := New(context.Background())
	if err != nil {
		t.Errorf("Cannot create client")
	}

	req := &runtimev0.CellServiceAllocateRequest{
		Cell: &runtimev0.Cell{
			Name:      testCellName,
			CpuQuota:  400000,
			CpuShares: 2,
		},
	}

	r, err := c.Runtime()
	if err != nil {
		t.Errorf("Runtime service not available")
	}

	rsp, err := r.AllocateCell(context.Background(), req)
	if err != nil {
		t.Errorf("AllocateCell should have NOT returned an error")
	}

	if rsp.CellName != testCellName {
		t.Error("Cell name wrong")
	}
}

func TestClientRuntimeFree(t *testing.T) {
	c, err := New(context.Background())
	if err != nil {
		t.Errorf("Cannot create client")
	}

	req := &runtimev0.CellServiceFreeRequest{
		CellName: testCellName,
	}

	r, err := c.Runtime()
	if err != nil {
		t.Errorf("Runtime service not available")
	}

	err = r.FreeCell(context.Background(), req)
	if err != nil {
		t.Errorf("FreeCell should have NOT returned an error")
	}
}
