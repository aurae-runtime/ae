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

package pki_create

import (
	"fmt"
	"io"

	aeCMD "github.com/aurae-runtime/ae/cmd"
	"github.com/aurae-runtime/ae/opt"
	"github.com/aurae-runtime/ae/pkg/pki"
	"github.com/spf13/cobra"
)

type outputVersion struct {
	BuildTime string `json:"buildTime,omitempty" yaml:"buildTime,omitempty"`
	Version   string `json:"version" yaml:"version"`
	Commit    string `json:"commit,omitempty" yaml:"commit,omitempty"`
}

type option struct {
	aeCMD.Option
	opt.OutputOption
	directory string
	domain    string
	writer    io.Writer
}

func (o *option) Complete(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("command 'create' requires a domain name as argument")
	}
	if len(args) > 1 {
		return fmt.Errorf("too many arguments for command 'create', expect %d, got %d", 1, len(args))
	}

	o.domain = args[0]

	return nil
}

func (o *option) Validate() error {
	return nil
}

func (o *option) Execute() error {
	pki.CreateRootCA(o.directory, o.domain)

	return nil
}

func (o *option) SetWriter(writer io.Writer) {
	o.writer = writer
}

func NewCMD() *cobra.Command {
	o := &option{}
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a CA for auraed.",
		Example: "ae create my.domain.com",

		RunE: func(cmd *cobra.Command, args []string) error {
			return aeCMD.Run(o, cmd, args)
		},
	}

	cmd.Flags().StringVarP(&o.directory, "dir", "d", o.directory, "Output directory to store CA files.")

	return cmd
}
