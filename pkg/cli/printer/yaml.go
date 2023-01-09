package printer

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

var _ Interface = NewYAML()

type YAML struct {
}

func NewYAML() *YAML {
	return &YAML{}
}

func (printer *YAML) Format() string {
	return "yaml"
}

func (printer *YAML) Print(obj any, w io.Writer) error {
	output, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(w, string(output))
	return err
}
