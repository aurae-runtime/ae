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
	discoveryv0 "github.com/aurae-runtime/ae/pkg/api/v0/discovery"
)

type outputHealthNode struct {
	Version string `json:"version"`
}

type outputHealth struct {
	Nodes	map[string]outputHealthNode `json:"nodes"`
}

func inc(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}

func hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	if len(ips) < 2 {
		return nil, fmt.Errorf("unexpectedly short ip list from cidr %q", cidr)
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

func protocol(ip_str string) (string, error) {
	ip := net.ParseIP(ip_str)
	if ip == nil {
		return "", fmt.Errorf("unable to parse IP %q", ip_str)
	}
	protocol := "tcp4";
	if ip.To4() == nil {
		protocol = "tcp6";
	}
	return protocol, nil
}

type option struct {
	aeCMD.Option
	cidr	string
	port	uint16
	verbose bool
	writer io.Writer
	outputFormat *cli.OutputFormat
}

func (o *option) Complete(_ []string) error {
	return nil
}

func (o *option) Validate() error {
	if err:= o.outputFormat.Validate(); err != nil {
		return err
	}

	if _, _, err := net.ParseCIDR(o.cidr); err != nil {
		return err
	}
	return nil
}

func (o *option) Execute() error {
	ctx := context.Background()
	hosts, err := hosts(o.cidr)
	if err != nil {
		return err
	}

	output := &outputHealth{
		Nodes: make(map[string]outputHealthNode),
	}

	// TODO: concurrency.
	for _, ip_str := range hosts {
		p, err := protocol(ip_str)
		if err != nil {
			return err
		}
		if (o.verbose) {
			log.Printf("connecting to %s:%d using protocol %s\n", ip_str, o.port, p)
		}
		c, err := client.New(ctx, config.WithSystem(config.System{Protocol: p, Socket: fmt.Sprintf("%s:%d", ip_str, o.port)}))
		if err != nil {
			return err
		}

		d, err := c.Discovery()
		if err != nil {
			return err
		}

		rsp, err := d.Health(ctx, &discoveryv0.HealthRequest{})
		if err != nil {
			if (o.verbose) {
				log.Printf("failed to call health: %s.  not an aurae node\n", err)
			}
			continue
		}

		if rsp.Healthy {
			output.Nodes[ip_str] = outputHealthNode{rsp.Version}
		}
	}

	return o.outputFormat.ToPrinter().Print(o.writer, output)
}

func (o *option) SetWriter(writer io.Writer) {
	o.writer = writer
}

func NewCMD() *cobra.Command {
	o := &option {
		outputFormat: cli.NewOutputFormat().
			WithDefaultFormat(printer.NewJSON().Format()).
			WithPrinter(printer.NewJSON()),
	}
	cmd := &cobra.Command{
		Use: "health",
		Short: "Scans a cluster of nodes for healthy Aurae nodes.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return aeCMD.Run(o, cmd, args)
		},
	}
	o.outputFormat.AddFlags(cmd)
	cmd.Flags().StringVar(&o.cidr, "cidr", o.cidr, "The network to scan in CIDR notation")
	cmd.Flags().Uint16Var(&o.port, "port", o.port, "The port to use when trying to connect")
	cmd.Flags().BoolVar(&o.verbose, "verbose", o.verbose, "Lots of output")
	return cmd
}
