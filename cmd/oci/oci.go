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

package oci

import (
	"fmt"
	"io"

	aeCMD "github.com/aurae-runtime/ae/cmd"
	"github.com/aurae-runtime/ae/cmd/oci/create"
	"github.com/aurae-runtime/ae/cmd/oci/delete"
	"github.com/aurae-runtime/ae/cmd/oci/kill"
	"github.com/aurae-runtime/ae/cmd/oci/start"
	"github.com/aurae-runtime/ae/cmd/oci/state"
	"github.com/spf13/cobra"
)

type option struct {
	aeCMD.Option
	writer io.Writer
}

func (o *option) Complete(_ []string) error {
	return nil
}

func (o *option) Validate() error {
	return nil
}

func (o *option) Execute() error {
	fmt.Fprintln(o.writer, "oci called")
	return nil
}

func (o *option) SetWriter(writer io.Writer) {
	o.writer = writer
}

func NewCMD() *cobra.Command {
	o := &option{}
	cmd := &cobra.Command{
		Use:   "oci",
		Short: "OCI Runtime Command Line Interface.",
		Long:  `OCI Runtime Command Line Interface.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return aeCMD.Run(o, cmd, args)
		},
	}

	cmd.AddCommand(create.NewCMD())
	cmd.AddCommand(delete.NewCMD())
	cmd.AddCommand(kill.NewCMD())
	cmd.AddCommand(start.NewCMD())
	cmd.AddCommand(state.NewCMD())

	return cmd
}
