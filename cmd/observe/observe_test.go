package observe

import (
	"testing"

	"github.com/aurae-runtime/ae/pkg/cli"
	"github.com/aurae-runtime/ae/pkg/cli/printer"
)

func TestComplete(t *testing.T) {
	ts := []struct {
		args        []string
		wantip      string
		wantlogtype string
		wanterr     bool
	}{
		{
			[]string{""},
			"", "",
			true,
		},
		{
			[]string{"foo"},
			"", "",
			true,
		},
		{
			[]string{"foo", "bar"},
			"", "",
			true,
		},
		{
			[]string{"foo", "bar", "baz"},
			"bar",
			"baz",
			false,
		},
		{
			[]string{"foo", "bar", "baz", "quux"},
			"", "",
			true,
		},
	}

	for _, tt := range ts {
		o := &option{}
		goterr := o.Complete(tt.args)
		if tt.wanterr && goterr == nil {
			t.Fatal("want error, got no error")
		}
		if !tt.wanterr && goterr != nil {
			t.Fatal("want no error, got error")
		}
		if tt.wantip != o.ip {
			t.Fatalf("want ip %q, got ip %q", tt.wantip, o.ip)
		}
		if tt.wantlogtype != o.logtype {
			t.Fatalf("want logtype %q, got logtype %q", tt.wantlogtype, o.logtype)
		}
	}
}

func TestValidate(t *testing.T) {
	ts := []struct {
		name         string
		outputFormat *cli.OutputFormat
		ip           string
		logtype      string
		wanterr      bool
	}{
		{
			name:         "no output format",
			outputFormat: cli.NewOutputFormat(),
			wanterr:      true,
		},
		{
			name:         "no ip or logtype",
			outputFormat: cli.NewOutputFormat().WithDefaultFormat("json").WithPrinter(printer.NewJSON()),
			wanterr:      true,
		},
		{
			name:         "invalid ip",
			outputFormat: cli.NewOutputFormat().WithDefaultFormat("json").WithPrinter(printer.NewJSON()),
			ip:           "invalid ip",
			wanterr:      true,
		},
		{
			name:         "invalid logtype",
			outputFormat: cli.NewOutputFormat().WithDefaultFormat("json").WithPrinter(printer.NewJSON()),
			ip:           "10.0.0.0",
			logtype:      "invalid",
			wanterr:      true,
		},
		{
			name:         "valid ip and logtype",
			outputFormat: cli.NewOutputFormat().WithDefaultFormat("json").WithPrinter(printer.NewJSON()),
			ip:           "10.0.0.0",
			logtype:      "daemon",
			wanterr:      false,
		},
	}

	for _, tt := range ts {
		o := &option{ip: tt.ip, logtype: tt.logtype, outputFormat: tt.outputFormat}
		goterr := o.Validate()
		if tt.wanterr && goterr == nil {
			t.Fatalf("[%s] want error, got no error", tt.name)
		}
		if !tt.wanterr && goterr != nil {
			t.Fatalf("[%s] want no error, got error: %s", tt.name, goterr)
		}
	}
}
