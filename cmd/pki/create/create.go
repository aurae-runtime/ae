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

package create

import (
	"context"
	"fmt"
	"io"

	aeCMD "github.com/aurae-runtime/ae/cmd"
	"github.com/aurae-runtime/ae/pkg/cli"
	"github.com/aurae-runtime/ae/pkg/cli/printer"
	"github.com/aurae-runtime/ae/pkg/pki"
	"github.com/spf13/cobra"
)

type option struct {
	aeCMD.Option
	outputFormat *cli.OutputFormat
	directory    string
	domain       string
	user         string
	silent       bool
	writer       io.Writer
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

func (o *option) Execute(_ context.Context) error {
	if o.user != "" {
		clientCSR, err := pki.HandleCreateClientCSR(o.directory, o.domain, o.user)
		if err != nil {
			return fmt.Errorf("failed to create client csr: %w", err)
		}
		if !o.silent {
			o.outputFormat.ToPrinter().Print(o.writer, &clientCSR)
		}

		return nil
	}

	rootCA, err := pki.HandleCreateAuraeRootCA(o.directory, o.domain)
	if err != nil {
		return fmt.Errorf("failed to create aurae root ca: %w", err)
	}
	if !o.silent {
		o.outputFormat.ToPrinter().Print(o.writer, &rootCA)
	}
	return nil
}

func (o *option) SetWriter(writer io.Writer) {
	o.writer = writer
}

func NewCMD(ctx context.Context) *cobra.Command {
	o := &option{
		outputFormat: cli.NewOutputFormat().
			WithDefaultFormat(printer.NewJSON().Format()).
			WithPrinter(printer.NewJSON()).
			WithPrinter(printer.NewYAML()),
		silent: false,
	}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a CA for auraed.",
		Example: `ae pki create my.domain.com
ae pki create --dir ./pki/ my.domain.com`,

		RunE: func(cmd *cobra.Command, args []string) error {
			return aeCMD.Run(ctx, o, cmd, args)
		},
	}

	o.outputFormat.AddFlags(cmd)
	cmd.Flags().StringVarP(&o.directory, "dir", "d", o.directory, "Output directory to store CA files.")
	cmd.Flags().StringVarP(&o.user, "user", "u", o.user, "Creates client certificate for a given user.")
	cmd.Flags().BoolVarP(&o.silent, "silent", "s", o.silent, "Silent mode, omits output")

	return cmd
}
