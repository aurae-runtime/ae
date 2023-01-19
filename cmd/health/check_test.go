package health

import (
	"testing"

	"github.com/aurae-runtime/ae/pkg/cli"
	"github.com/aurae-runtime/ae/pkg/cli/printer"
)

func TestComplete(t *testing.T) {
	ts := []struct {
		args    []string
		wanterr bool
	}{
		{
			args:    []string{"cidr", "192.168.0.0/32", "list,of,services"},
			wanterr: false,
		},
		{
			args:    []string{"ip", "10.0.0.0", "list,of,services"},
			wanterr: false,
		},
		{
			args:    []string{"cidr"},
			wanterr: true,
		},
		{
			args:    []string{"cidr", "192.168.0.0/32"},
			wanterr: true,
		},
	}

	for _, tt := range ts {
		o := &option{}
		goterr := o.Complete(tt.args)
		if tt.wanterr && goterr == nil {
			t.Fatal("want error, got no error")
		}
		if !tt.wanterr && goterr != nil {
			t.Fatalf("want no error, got error %q", goterr)
		}
	}
}

func TestValidate(t *testing.T) {
	ts := []struct {
		name         string
		outputFormat *cli.OutputFormat
		cidr         string
		ip           string
		services     []string
		wanterr      bool
	}{
		{
			name:         "no output format",
			outputFormat: cli.NewOutputFormat(),
			wanterr:      true,
		},
		{
			name:         "no cidr or ip",
			outputFormat: cli.NewOutputFormat().WithDefaultFormat("json").WithPrinter(printer.NewJSON()),
			wanterr:      true,
		},
		{
			name:         "no services",
			outputFormat: cli.NewOutputFormat().WithDefaultFormat("json").WithPrinter(printer.NewJSON()),
			cidr:         "192.168.170.0/32",
			wanterr:      true,
		},
		{
			name:         "invalid cidr",
			outputFormat: cli.NewOutputFormat().WithDefaultFormat("json").WithPrinter(printer.NewJSON()),
			cidr:         "invalid cidr",
			wanterr:      true,
		},
		{
			name:         "valid cidr",
			outputFormat: cli.NewOutputFormat().WithDefaultFormat("json").WithPrinter(printer.NewJSON()),
			cidr:         "192.168.170.0/32",
			services:     []string{"foo", "bar"},
			wanterr:      false,
		},
		{
			name:         "invalid ip",
			outputFormat: cli.NewOutputFormat().WithDefaultFormat("json").WithPrinter(printer.NewJSON()),
			ip:           "invalid ip",
			wanterr:      true,
		},
		{
			name:         "valid ip",
			outputFormat: cli.NewOutputFormat().WithDefaultFormat("json").WithPrinter(printer.NewJSON()),
			ip:           "10.0.0.0",
			services:     []string{"foo", "bar"},
			wanterr:      false,
		},
	}

	for _, tt := range ts {
		o := &option{cidr: tt.cidr, ip: tt.ip, outputFormat: tt.outputFormat, services: tt.services}
		goterr := o.Validate()
		if tt.wanterr && goterr == nil {
			t.Fatalf("[%s] want error, got no error", tt.name)
		}
		if !tt.wanterr && goterr != nil {
			t.Fatalf("[%s] want no error, got error: %s", tt.name, goterr)
		}
	}
}
