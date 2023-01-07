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
	"io"

	aeCMD "github.com/aurae-runtime/ae/cmd"
	ociCreate "github.com/aurae-runtime/ae/cmd/oci/create"
	ociDelete "github.com/aurae-runtime/ae/cmd/oci/delete"
	ociKill "github.com/aurae-runtime/ae/cmd/oci/kill"
	ociStart "github.com/aurae-runtime/ae/cmd/oci/start"
	ociState "github.com/aurae-runtime/ae/cmd/oci/state"
	"github.com/aurae-runtime/ae/opt"
	"github.com/aurae-runtime/ae/output"
	"github.com/spf13/cobra"
)

type option struct {
	aeCMD.Option
	opt.OutputOption
	writer io.Writer
}

func (o *option) Complete(_ []string) error {
	return nil
}

func (o *option) Validate() error {
	return nil
}

func (o *option) Execute() error {
	return output.HandleString(o.writer, "oci called")
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
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cmd.AddCommand(ociCreate.NewCMD())
	cmd.AddCommand(ociDelete.NewCMD())
	cmd.AddCommand(ociKill.NewCMD())
	cmd.AddCommand(ociStart.NewCMD())
	cmd.AddCommand(ociState.NewCMD())

	return cmd
}
