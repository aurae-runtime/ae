package cli

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/aurae-runtime/ae/pkg/cli/printer"
	"github.com/spf13/cobra"
)

func TestOutputFormatAddFlags(t *testing.T) {
	tests := []struct {
		name        string
		printers    []printer.Interface
		input       []string
		expectedOut string
	}{
		{
			name:     "'output' is valid flag",
			printers: []printer.Interface{printer.NewJSON()},
			input:    []string{"--output", "json"},
		},
		{
			name:     "'o' is valid flag",
			printers: []printer.Interface{printer.NewJSON()},
			input:    []string{"-o", "json"},
		},
	}

	for _, test := range tests {
		cmd := &cobra.Command{}
		output := NewOutputFormat()
		for _, p := range test.printers {
			output.WithPrinter(p)
		}
		output.AddFlags(cmd)

		buf := &bytes.Buffer{}
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs(test.input)

		if err :=cmd.Execute(); err != nil {
			t.Errorf("%s: errored. Got: %s", test.name, buf.String())
		}

		if buf.String() != "" {
			t.Errorf("%s: expected no output. Got: %s", test.name, buf.String())
		}
	}
}

func TestOutputFormatValidate(t *testing.T) {
	tests := []struct {
		format      string
		printers    []printer.Interface
		expectedErr error
	}{
		{
			format:      "json",
			printers:    []printer.Interface{printer.NewJSON(), printer.NewYAML()},
			expectedErr: nil,
		},
		{
			format:      "yaml",
			printers:    []printer.Interface{printer.NewJSON(), printer.NewYAML()},
			expectedErr: nil,
		},
		{
			format:      "json",
			printers:    []printer.Interface{printer.NewYAML()},
			expectedErr: errors.New("\"json\": unknown output format"),
		},
	}

	for _, test := range tests {
		output := NewOutputFormat()
		for _, p := range test.printers {
			output.WithPrinter(p)
		}
		output.format = test.format

		err := output.Validate()
		if test.expectedErr == nil {
			if err != nil {
				t.Errorf("expected: nil; got: %v", err)
			}
		} else if test.expectedErr.Error() != err.Error() {
			t.Errorf("expected: %v; got: %v", test.expectedErr, err)
		}
	}
}

func TestOutputFormatToPrinter(t *testing.T) {
	tests := []struct {
		format       string
		printers     []printer.Interface
		expectedType reflect.Type
	}{
		{
			format:       "json",
			printers:     []printer.Interface{printer.NewJSON(), printer.NewYAML()},
			expectedType: reflect.TypeOf(printer.JSON{}),
		},
		{
			format:       "yaml",
			printers:     []printer.Interface{printer.NewJSON(), printer.NewYAML()},
			expectedType: reflect.TypeOf(printer.YAML{}),
		},
	}

	for _, test := range tests {
		output := NewOutputFormat()
		for _, p := range test.printers {
			output.WithPrinter(p)
		}
		output.format = test.format

		printerType := reflect.Indirect(reflect.ValueOf(output.ToPrinter())).Type()
		if test.expectedType != printerType {
			t.Errorf("expected: %v; got: %v", test.expectedType, printerType)
		}
	}
}

func TestOutputFormatFlagHelp(t *testing.T) {
	cmd := &cobra.Command{}
	output := NewOutputFormat().
		WithPrinter(printer.NewJSON()).
		WithPrinter(printer.NewYAML())
	output.AddFlags(cmd)
	f := cmd.Flag("output")
	expected := "Output format. One of: (json, yaml)."
	if f.Usage != expected {
		t.Errorf("expected: %s; got: %s", expected, f.Usage)
	}
}
