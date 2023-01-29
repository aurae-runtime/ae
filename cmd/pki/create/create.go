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
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"

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
	caPath       string
	caKeyPath    string
	csrPath      string
	csr          *pki.CertificateRequest
	ca           *pki.Certificate
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

	if o.caPath != "" && (o.caKeyPath == "" || o.csrPath == "") {
		return fmt.Errorf("must provide --caKey and --csr when using --ca")
	}

	if o.caPath != "" {
		b, err := os.ReadFile(o.caPath)
		if err != nil {
			return fmt.Errorf("failed to read ca certificate: %w", err)
		}

		o.ca = &pki.Certificate{}
		o.ca.Certificate = string(b)
	}

	if o.caKeyPath != "" {
		b, err := os.ReadFile(o.caKeyPath)
		if err != nil {
			return fmt.Errorf("failed to read ca private key: %w", err)
		}

		o.ca.PrivateKey = string(b)
	}

	if o.csrPath != "" {
		b, err := os.ReadFile(o.csrPath)
		if err != nil {
			return fmt.Errorf("failed to read csr: %w", err)
		}

		o.csr = &pki.CertificateRequest{}
		o.csr.CSR = string(b)
	}

	o.domain = args[0]

	return nil
}

func (o *option) Validate() error {
	if o.caPath != "" {
		caPem, _ := pem.Decode([]byte(o.ca.Certificate))
		_, err := x509.ParseCertificate(caPem.Bytes)
		if err != nil {
			return fmt.Errorf("could not parse ca file")
		}
	}

	if o.caKeyPath != "" {
		caKeyPem, _ := pem.Decode([]byte(o.ca.PrivateKey))
		_, err := x509.ParsePKCS1PrivateKey(caKeyPem.Bytes)
		if err != nil {
			return fmt.Errorf("could not parse ca file")
		}
	}

	if o.csrPath != "" {
		csrPem, _ := pem.Decode([]byte(o.csr.CSR))
		_, err := x509.ParseCertificateRequest(csrPem.Bytes)
		if err != nil {
			return fmt.Errorf("could not parse csr file")
		}
	}

	return nil
}

func (o *option) Execute(_ context.Context) error {
	if o.user != "" {

		if o.caPath != "" {
			clientCrt, err := pki.CreateClientCertificate(o.directory, o.csr.CSR, o.ca, o.user)
			if err != nil {
				return fmt.Errorf("failed to create client certificate: %w", err)
			}
			if !o.silent {
				o.outputFormat.ToPrinter().Print(o.writer, &clientCrt)
			}

			return nil
		}

		clientCSR, err := pki.CreateClientCSR(o.directory, o.domain, o.user)
		if err != nil {
			return fmt.Errorf("failed to create client csr: %w", err)
		}
		if !o.silent {
			o.outputFormat.ToPrinter().Print(o.writer, &clientCSR)
		}

		return nil
	}

	rootCA, err := pki.CreateAuraeRootCA(o.directory, o.domain)
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
	cmd.Flags().StringVar(&o.caPath, "ca", o.caPath, "Use the given CA certificate.")
	cmd.Flags().StringVar(&o.caKeyPath, "caKey", o.caKeyPath, "The corresponding CA key.")
	cmd.Flags().StringVar(&o.csrPath, "csr", o.csrPath, "CSR input file.")
	cmd.Flags().BoolVarP(&o.silent, "silent", "s", o.silent, "Silent mode, omits output")

	return cmd
}
