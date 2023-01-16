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

package testsuite

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

type Test struct {
	Title          string
	Cmd            *cobra.Command
	Args           []string
	ExpectedStdout string
}

func ExecuteSuite(t *testing.T, tests []Test) {
	for _, test := range tests {
		t.Run(test.Title, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			test.Cmd.SetOut(buffer)
			test.Cmd.SetArgs(test.Args)

			err := test.Cmd.Execute()

			if !IsNil(err) {
				t.Errorf("Unexpected error while executing command: %v", err)
			}
			if !IsEqualString(test.ExpectedStdout, buffer.String()) {
				t.Errorf("Unexpected Stdout\ngot:\n%s\nwant:\n%s\n", test.ExpectedStdout, buffer.String())
			}
		})
	}
}

func IsEqualString(s1, s2 string) bool {
	return s1 == s2
}

func IsNil(object any) bool {
	return object == nil
}
