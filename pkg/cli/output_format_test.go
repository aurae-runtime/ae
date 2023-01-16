package cli

import (
	"errors"
	"reflect"
	"testing"

	"github.com/aurae-runtime/ae/pkg/cli/printer"
	"github.com/aurae-runtime/ae/pkg/cli/testsuite"
	"github.com/spf13/cobra"
)

func newCmd(printers []printer.Interface) *cobra.Command {
	cmd := &cobra.Command{}
	output := NewOutputFormat()
	for _, p := range printers {
		output.WithPrinter(p)
	}
	output.AddFlags(cmd)
	return cmd
}

func TestOutputFormatAddFlags(t *testing.T) {
	tests := []testsuite.Test{
		{
			Title:          "'output' is valid flag",
			Cmd:            newCmd([]printer.Interface{printer.NewJSON()}),
			Args:           []string{"--output", "json"},
			ExpectedStdout: "",
		},
		{
			Title:          "'o' is valid flag",
			Cmd:            newCmd([]printer.Interface{printer.NewJSON()}),
			Args:           []string{"-o", "json"},
			ExpectedStdout: "",
		},
	}

	testsuite.ExecuteSuite(t, tests)
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
	tests := []testsuite.Test{
		{
			Title:          "Help is formatted correctly",
			Cmd:            newCmd([]printer.Interface{printer.NewJSON(), printer.NewYAML()}),
			Args:           []string{"help"},
			ExpectedStdout: "",
		},
	}

	testsuite.ExecuteSuite(t, tests)
}
