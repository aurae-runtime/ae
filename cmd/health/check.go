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

package health

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/3th1nk/cidr"
	"github.com/aurae-runtime/ae/pkg/cli"
	"github.com/aurae-runtime/ae/pkg/cli/printer"
	"github.com/aurae-runtime/ae/pkg/client"
	"github.com/aurae-runtime/ae/pkg/config"
	"github.com/spf13/cobra"

	aeCMD "github.com/aurae-runtime/ae/cmd"
	//healthv1 "github.com/aurae-runtime/ae/pkg/api/grpc/health/v1/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

type outputCheckNode struct {
	// Maps from a service name to a serving status.
	Statuses map[string]string `json:"version"`
}

type outputCheck struct {
	Nodes map[string]outputCheckNode `json:"nodes"`
}

type option struct {
	aeCMD.Option

	// TODO: abstract the next batch of fields into a "clusterOption"
	ctx          context.Context
	cidr         string
	ip           string
	port         uint16
	protocol     string
	verbose      bool
	writer       io.Writer
	outputFormat *cli.OutputFormat

	services []string
	output   *outputCheck
}

func (o *option) Complete(args []string) error {
	if len(args) != 3 {
		return errors.New("expected 'cidr' or 'ip' and a list of services passed to this command")
	}
	switch args[0] {
	case "cidr":
		o.cidr = args[1]
	case "ip":
		o.ip = args[1]
	default:
		return errors.New("either 'cidr' or 'ip' must be passed to this command")
	}

	o.services = strings.Split(args[2], ",")

	return nil
}

func (o *option) Validate() error {
	if err := o.outputFormat.Validate(); err != nil {
		return err
	}

	if len(o.cidr) != 0 {
		if _, _, err := net.ParseCIDR(o.cidr); err != nil {
			return err
		}
	} else if len(o.ip) != 0 {
		if ip := net.ParseIP(o.ip); ip == nil {
			return fmt.Errorf("failed to parse ip %q", o.ip)
		}
	} else {
		return errors.New("either 'cidr' or 'ip' must be passed to this command")
	}

	// TODO: check for valid service list.
	if len(o.services) == 0 {
		return errors.New("expected list of services to be provided")
	}

	return nil
}

func (o *option) Execute(ctx context.Context) error {
	o.ctx = ctx

	o.output = &outputCheck{
		Nodes: make(map[string]outputCheckNode),
	}

	o.protocol = "tcp4"

	if len(o.cidr) != 0 {
		c, err := cidr.Parse(o.cidr)
		if err != nil {
			return err
		}
		if c.IsIPv6() {
			o.protocol = "tcp6"
		}
		// TODO: use some concurrency magic to run these async
		c.Each(o.checkHost)
	} else if len(o.ip) != 0 {
		if net.ParseIP(o.ip).To4() == nil {
			o.protocol = "tcp6"
		}
		o.checkHost(o.ip)
	}

	return o.outputFormat.ToPrinter().Print(o.writer, o.output)
}

func (o *option) SetWriter(writer io.Writer) {
	o.writer = writer
}

func (o option) checkHost(ip_str string) bool {
	if o.verbose {
		log.Printf("connecting to %s:%d using protocol %s\n", ip_str, o.port, o.protocol)
	}

	c, err := client.New(o.ctx, config.WithSystem(config.System{Protocol: o.protocol, Socket: net.JoinHostPort(ip_str, fmt.Sprintf("%d", o.port))}))
	if err != nil {
		log.Fatalf("failed to create client: %s", err)
		return false
	}

	h, err := c.Health()
	if err != nil {
		log.Fatalf("failed to dial Health service: %s", err)
		return false
	}

	o.output.Nodes[ip_str] = outputCheckNode{
		Statuses: make(map[string]string),
	}

	for _, s := range o.services {
		rsp, err := h.Check(o.ctx, &healthv1.HealthCheckRequest{Service: s})
		if err != nil {
			if o.verbose {
				log.Printf("failed to call discover: %s.  not an aurae node\n", err)
			}
			return true
		}

		o.output.Nodes[ip_str].Statuses[s] = rsp.Status.String()
	}
	return true
}

func NewCMD(ctx context.Context) *cobra.Command {
	o := &option{
		outputFormat: cli.NewOutputFormat().
			WithDefaultFormat(printer.NewJSON().Format()).
			WithPrinter(printer.NewJSON()),
	}
	cmd := &cobra.Command{
		Use:   "check [cidr <cidr>|ip <ip>] [services]",
		Short: "Scans a node or cluster of nodes and checks the health of the given list of services",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return aeCMD.Run(ctx, o, cmd, args)
		},
	}
	o.outputFormat.AddFlags(cmd)
	cmd.Flags().Uint16Var(&o.port, "port", o.port, "The port to use when trying to connect")
	cmd.Flags().BoolVar(&o.verbose, "verbose", o.verbose, "Lots of output")
	return cmd
}
