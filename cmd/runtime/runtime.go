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
	"io"

	aeCMD "github.com/aurae-runtime/ae/cmd"
	allocate "github.com/aurae-runtime/ae/cmd/runtime/allocate"
	free "github.com/aurae-runtime/ae/cmd/runtime/free"
	start "github.com/aurae-runtime/ae/cmd/runtime/start"
	stop "github.com/aurae-runtime/ae/cmd/runtime/stop"
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
	return output.HandleString(o.writer, "runtime called")
}

func (o *option) SetWriter(writer io.Writer) {
	o.writer = writer
}

func NewCMD() *cobra.Command {
	o := &option{}
	cmd := &cobra.Command{
		Use:   "runtime",
		Short: "CLI for aurae runtime services.",
		Long:  `CLI for aurae runtime services.`,
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
	cmd.AddCommand(allocate.NewCMD())
	cmd.AddCommand(free.NewCMD())
	cmd.AddCommand(start.NewCMD())
	cmd.AddCommand(stop.NewCMD())

	return cmd
}
