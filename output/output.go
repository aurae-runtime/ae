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

package output

import (
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

const (
	JSONOutput  = "json"
	YAMLOutput  = "yaml"
	TableOutput = "table"
	GoOutput    = "go"
)

// ValidateAndSet will validate the given output and if it's empty will set it
// with the default value "yaml"
func ValidateAndSet(o *string) error {
	if *o == "" {
		*o = YAMLOutput
		return nil
	} else if *o != YAMLOutput &&
		*o != JSONOutput &&
		*o != TableOutput &&
		*o != GoOutput {
		return fmt.Errorf(
			"--ouput must be %q, %q, %q or %q",
			JSONOutput,
			YAMLOutput,
			TableOutput,
			GoOutput,
		)
	}

	return nil
}

func Handle(writer io.Writer, output string, obj interface{}) error {
	var data []byte
	var err error

	switch output {
	case JSONOutput:
		data, err = json.Marshal(obj)
	case TableOutput:
		return fmt.Errorf("table output not implemented yet")
	case GoOutput:
		return fmt.Errorf("go output not implemented yet")
	default:
		data, err = yaml.Marshal(obj)
	}
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, string(data))
	return err
}

func HandleString(writer io.Writer, msg string) error {
	_, err := fmt.Fprintln(writer, msg)
	return err
}
