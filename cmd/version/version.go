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

package version

import (
	"context"
	"io"

	aeCMD "github.com/aurae-runtime/ae/cmd"
	"github.com/aurae-runtime/ae/pkg/cli"
	"github.com/aurae-runtime/ae/pkg/cli/printer"
	"github.com/prometheus/common/version"
	"github.com/spf13/cobra"
)

type outputVersion struct {
	BuildTime string `json:"buildTime,omitempty" yaml:"buildTime,omitempty"`
	Version   string `json:"version" yaml:"version"`
	Commit    string `json:"commit,omitempty" yaml:"commit,omitempty"`
}

type option struct {
	aeCMD.Option
	writer       io.Writer
	short        bool
	outputFormat *cli.OutputFormat
}

func (o *option) Complete(_ []string) error {
	return nil
}

func (o *option) Validate() error {
	if err := o.outputFormat.Validate(); err != nil {
		return err
	}
	return nil
}

func (o *option) Execute(_ context.Context) error {
	clientVersion := &outputVersion{
		Version: version.Version,
	}
	if !o.short {
		clientVersion.BuildTime = version.BuildDate
		clientVersion.Commit = version.Revision
	}

	return o.outputFormat.ToPrinter().Print(o.writer, clientVersion)
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
	}
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Displays Aurae command line client version.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return aeCMD.Run(ctx, o, cmd, args)
		},
	}
	o.outputFormat.AddFlags(cmd)
	cmd.Flags().BoolVar(&o.short, "short", o.short, "If true, just print the version number.")
	return cmd
}
