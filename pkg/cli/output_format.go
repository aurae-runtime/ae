package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aurae-runtime/ae/pkg/cli/printer"
	"github.com/spf13/cobra"
)

var ErrUnknownFormat = errors.New("unknown output format")

type OutputFormat struct {
	format   string
	printers []printer.Interface
}

func NewOutputFormat() *OutputFormat {
	return &OutputFormat{
		"",
		nil,
	}
}

func (o *OutputFormat) WithDefaultFormat(defaultFormat string) *OutputFormat {
	o.format = defaultFormat
	return o
}

func (o *OutputFormat) WithPrinter(p printer.Interface) *OutputFormat {
	o.printers = append(o.printers, p)
	return o
}

func (o *OutputFormat) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(
		&o.format,
		"output",
		"o",
		o.format,
		fmt.Sprintf("Output format. One of: (%s).", strings.Join(o.allowedFormats(), ", ")),
	)
}

func (o *OutputFormat) Validate() error {
	if o.format == "" {
		return nil
	}

	for _, printer := range o.printers {
		if o.format == printer.Format() {
			return nil
		}
	}

	return fmt.Errorf("%q: %w", o.format, ErrUnknownFormat)
}

func (o *OutputFormat) ToPrinter() printer.Interface {
	for _, printer := range o.printers {
		if o.format == printer.Format() {
			return printer
		}
	}
	return nil
}

func (o *OutputFormat) allowedFormats() []string {
	var allowed []string
	for _, printer := range o.printers {
		allowed = append(allowed, printer.Format())
	}
	return allowed
}
