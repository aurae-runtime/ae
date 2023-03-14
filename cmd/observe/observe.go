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

package observe

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/aurae-runtime/ae/pkg/cli"
	"github.com/aurae-runtime/ae/pkg/cli/printer"
	"github.com/aurae-runtime/ae/pkg/client"
	"github.com/aurae-runtime/ae/pkg/config"
	"github.com/spf13/cobra"

	aeCMD "github.com/aurae-runtime/ae/cmd"
	observev0 "github.com/aurae-runtime/ae/pkg/api/v0/observe"
)

type outputObserve struct{}

type option struct {
	aeCMD.Option
	ctx          context.Context
	ip           string
	logtype      string
	port         uint16
	protocol     string
	verbose      bool
	writer       io.Writer
	output       *outputObserve
	outputFormat *cli.OutputFormat
}

func (o *option) Complete(args []string) error {
	log.Println(args)
	if len(args) != 2 {
		return errors.New("expected ip address and log type to be passed to this command")
	}
	o.ip = args[0]
	o.logtype = args[1]
	return nil
}

func (o *option) Validate() error {
	if err := o.outputFormat.Validate(); err != nil {
		return err
	}
	if len(o.ip) == 0 {
		return errors.New("ip address must be passed to this command")
	}
	if ip := net.ParseIP(o.ip); ip == nil {
		return fmt.Errorf("failed to parse IP %q", o.ip)
	}
	if o.logtype != "daemon" && o.logtype != "subprocesses" {
		return errors.New("either 'daemon' or 'subprocesses' must be passed to the command")
	}
	return nil
}

func (o *option) Execute(ctx context.Context) error {
	o.ctx = ctx
	o.output = &outputObserve{}

	o.protocol = "tcp4"
	if net.ParseIP(o.ip).To4() == nil {
		o.protocol = "tcp6"
	}
	o.observeHost(o.ip)
	return o.outputFormat.ToPrinter().Print(o.writer, o.output)
}

func (o *option) SetWriter(writer io.Writer) {
	o.writer = writer
}

func (o option) observeHost(ip_str string) bool {
	if o.verbose {
		log.Printf("connecting to %s:%d using protocol %s\n", ip_str, o.port, o.protocol)
	}

	c, err := client.New(o.ctx, config.WithSystem(config.System{Protocol: o.protocol, Socket: net.JoinHostPort(ip_str, fmt.Sprintf("%d", o.port))}))
	if err != nil {
		log.Fatalf("failed to create client: %s", err)
		return false
	}

	obs, err := c.Observe()
	if err != nil {
		log.Fatalf("failed to dial Observe service: %s", err)
		return false
	}

	// TODO: handle output format
	switch o.logtype {
	case "daemon":
		// TODO: request parameters
		req := observev0.GetAuraeDaemonLogStreamRequest{} // TODO
		stream, err := obs.GetAuraeDaemonLogStream(o.ctx, &req)
		if err != nil {
			log.Fatalf("%s", err)
			return false
		}
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("%s", err)
				return false
			}
			if _, err := o.writer.Write([]byte(resp.Item.Line)); err != nil {
				log.Fatalf("%s", err)
				return false
			}
		}
	case "subprocesses":
		// TODO: request parameters
		req := observev0.GetSubProcessStreamRequest{} // TODO
		stream, err := obs.GetSubProcessStream(o.ctx, &req)
		if err != nil {
			log.Fatalf("%s", err)
			return false
		}
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("%s", err)
				return false
			}
			if _, err := o.writer.Write([]byte(resp.Item.Line)); err != nil {
				log.Fatalf("%s", err)
				return false
			}
		}
	}

	return true
}

func NewCMD(ctx context.Context) *cobra.Command {
	o := &option{
		outputFormat: cli.NewOutputFormat().WithDefaultFormat(printer.NewJSON().Format()).WithPrinter(printer.NewJSON()),
	}
	cmd := &cobra.Command{
		Use:   "observe <ip> <daemon|subprocesses>",
		Short: "get a stream of logs either from the aurae daemon or spawned subprocesses running on the given IP address",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return aeCMD.Run(ctx, o, cmd, args)
		},
	}
	o.outputFormat.AddFlags(cmd)
	cmd.Flags().Uint16Var(&o.port, "port", o.port, "The port to use when connecting")
	cmd.Flags().BoolVar(&o.verbose, "verbose", o.verbose, "Lots of output")
	return cmd
}
